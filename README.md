Note: This is a school project completed in Summer of 2022. This project was completed in collaboration with Anson Quon at UC Berkeley.

The starter code that was provided with this project is located at https://github.com/cs161-staff/project2-starter-code




For comprehensive documentation, see the Project 2 Spec (https://cs161.org/proj2/getting-started-coding/).

To test the implementation, run `go test -v` inside of the `client_test` directory.

Implementation-Specific Documentation:

Data Structures:
-User: (should be stored encrypted using a key generated with the password. Password is not stored anywhere. UUID should be generated with the username.) 

Attributes:

-Username

-user’s signing key SigningKey

-user’s authPointer symmetric key FileKey -user’s invitation private key. InvitationKey

-FileChunk Attributes:

-NextChunkUUID (UUID corresponding to the next file chunk stored in DataStore) -Content

-LastChunk UUID (UUID corresponding to the last chunk of the entire File)

-FileAuth Attributes:

-FileUUID (UUID of the file this “middleman” object points to)

-Filekey to access file

-Hash map of usernames to TopAuthDetails objects -IsOwner boolean flag

-AuthPointer Attributes:

-AuthUUID (UUID of the FileAuth “middleman” this AuthPointer points to) 

-Authkey for encrypting/decrypting the FileAuth ciphertext.

-Invitation: Attributes

-AuthUUID (UUID corresponding to the Auth Pointer)

-Authkey for encrypting/decrypting Auth Pointer ciphertext. -TopAuthDetails (abstraction layer holding the UUID and key for a user’s FileAuth):

Attributes

-AuthUUID for FileAuth
-AuthKey for FileAuth
“Mailbox or Middleman” object fileAuth, sharing and revocation:
- Generated on a per-file and per top level user basis, the owner and each top level authorized user has access to a different fileAuth.
- When the owner shares the file, they create a new fileAuth for the new top level user.
- When any other user shares the file, they give the new user the location and the fileAuth private key that
belongs to the top level authorized user.
- Accepting an invitation will create a new authpointer object for the new user.
- Owner’s fileAuth contains a hashmap of top-level authorized users to TopAuthDetails containing the
corresponding FileAuth UUID and the key to decrypt the encrypted FileAuth.
- fileAuths are encrypted using the fileAuth private key shared with the top level user.
- All user’s fileAuth contains the file UUID, file access key, a map of user uuids to TopAuthDetails, and an
isOwner boolean flag.
- If a user is revoked, owner will reencrypt and move the file, then create new fileAuth objects containing the
the updated file location and key, then encrypt that fileAuth object using the same shared private keys used
to encrypt the original and store them at the same location
- All users that had access to the prior fileAuth still maintain access to the new one with its updated fields.
Per user authPointer:
-located at UUID(H(filename||authorizedusername||”ptr”))
-created for each user a file is shared with, generated upon accepting an invitation. -contains the fileAuth private key and the location of the fileAuth object.
(see fig.1) for relations between users, authPointers, fileAuths and files.
Helper methods (all store methods implement integrity and confidentiality via encrypting the struct. All get methods verify integrity and decrypt the struct):
 -createFileAuth(file uuid, file access key, map[string]TopAuthDetails,isOwner) -> returns fileAuth -storeFileAuth(uuid, FileAuth private key, FileAuth object) -> returns error
-getFileAuth(uuid, FileAuth private key) -> returns (FileAuth, error)
-createAuthPointer(fileauth UUID, fileAuth private key) -> returns AuthPointer -storeAuthPointer(username string, filename string, AuthPointer symmetric key, authPointer object) -> returns error
-getAuthPointer(username string, filename string, AuthPointer symmetric key) -> returns (AuthPointer, exists bool, error)
-createFileChunk(nextuuid, content, lastuuid) -> returns FileChunk
-storeFileChunk(uuid, key, chunk) -> returns error
-getFileChunk(uuid, key) -> returns (FileChunk, exists bool, error) -createTopAuthDetails(fileAuthUUID, fileAuthKey) -> returns TopAuthDetails -createInvitationObject(fileAuthUUID, fileAuthKey) -> returns Invitation -storeInvitationObject(uuid, publicKey, invitation) -> returns error -getInvitationObject(uuid, dkey, username) -> returns (Invitation, error)
Authentication.
(pointers are just UUIDs)
InitUser(string username, string password) -> returns user pointer, error message
1. Check if length of username is 0; if so, return an error
2. Check if a user with the same username exists; if so, return an error.
3. Check that the length of the username is 0; if so, return an error.
4. Creates user object with all necessary fields
5. Store the user object’s signing and file access public key in KeyStore.
6. Generates a UUID using the username
7. Encrypted the user object using symmetric key encryption with a key generated deterministically with the password
6. Store the user object in datastore with UUID, encrypted user object. MAC the user with HashKDF of userEncryptionKey. The resulting key for the MAC is different from the key user to encrypt the User object.
GetUser(string username, string password), returns user pointer and error message.
1. Looks for the UUID generated with the username in datastore
2. If not found, return error
3. Fetch encrypted user object from datastore.
4. Check for the integrity of the MACed user object using the derived MAC key.. Error if there is a loss of integrity (Verification)
5. Attempt to decrypt the user object fetched from datastore using the provided password
6. If the username field of the object does not match the username (file decrypted with wrong key, would be
scrambled or unreadable), return error
7. If the username field of the object matches up with the provided username, success.
Appending:
Have the file stored in chunks separated into fixed length nodes at random UUIDs like a linkedlist, with fileChunks containing info about the next node and the last node in the case of the first chunk.
When calling append, decrypt the first node of the file, read the location of lastNode (node after the current last node of the file), store the appended data to that address, then update the last node location to the new address.
File Storage:
Files are stored in chunks, in the format of a linked list. Each file chunk contains 512 bytes of content and the location of the next chunk, and the first file chunk of a file contains the location of the next chunk to come after the last chunk. A loop appends content starting from the file’s starting point for a StoreFile(), and starting from the chunk after the last already stored chunk for an AppendTo().
<img width="653" alt="image" src="https://user-images.githubusercontent.com/58440224/206352082-96c8889f-e00a-4a76-b030-3d438c2e5209.png">

(fig.1, relations between users, authPointers, fileAuths and files.)

Test Proposal:

 Design Requirement: The client MUST enforce that there is only a single copy of a file. Sharing the file MAY NOT create a copy of the file. [3.6.4]
1. Create 2 users
2. Have one user share a file with the other user.
3. Have one user modify the shared file
4. Check that the other user can see the changes on the shared file.
5. If the other user can see the changes, then both users see the same copy.
Design Requirement: The client MUST ensure confidentiality and integrity of file contents and file sharing invitations. [3.5.1]
1. Create a user
2. Create a file
3. Create another user, and share the file to that user (send an invitation)
4. Expect that the entries stored in the locations of the file and the invitation are scrambled when downloaded
and do not leak any information. In an ideal scenario, the adversary would not even have access to the locations.
Design Requirement: The client must allow users to efficiently append new content to previously stored files. [3.7.1]
1. Create a user
2. Create a file of length 10 bytes from the user
3. Append 10 bytes to that file
4. Measure the bandwidth used
5. Append 1,000,000 bytes to the file
6. Append 10 bytes to that file
7. Measure the bandwidth used by the latest append.
8. Expect that there is no significant difference in bandwidth used between the 1st 10 bytes append and the
2nd 10 byte append.
Design Requirement: The client MAY leak any information except filenames, lengths of filenames, file contents, and file sharing invitations. For example, the client design MAY leak the size of file contents or the number of files associated with a user. [3.5.5]
1. Create multiple file objects of different content and length.
2. Encrypt the file objects using our chosen encryption algorithms
3. Check that the lengths of each of the encrypted files are the same.
4. Possibly perform some attacks like length-extension to check for robustness
Design Requirement: The client MUST NOT save any data to the local file system. If the client is restarted, it must be able to pick up where it left off given only a username and password.
1. Create a user
2. Have the user create a file with some content
3. Append additional information to the file
4. Log out of the client
5. Log back into the the user
6. Download the file
7. Expect all contents to be present in the file.
Design Requirement: The client MUST ensure the integrity of filenames [3.5.2]

1. Have a file object created by a fileauth, and have the fileauth sign the encrypted file object with the fileauth’s private key.
2. Tamper with the encrypted file object.
3. When the fileauth retrieves the file and signature, verify the signature using the fileauth’s public key.
