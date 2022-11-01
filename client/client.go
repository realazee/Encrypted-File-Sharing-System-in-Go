package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username      string
	SigningKey    userlib.DSSignKey
	FileKey       []byte //accesses the AuthPointer object
	InvitationKey userlib.PKEDecKey

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

//Created whenever a File is created. A FileAuth is generated on a per-file and per top level user basis; an
//owner and all their top-level authorized users have access to one FileAuth, and this structure is repeated for
//users who share the file with other users.
//FileAuths will be stored at random locations and encrypted with a public key, both shared as a part of an invitation.
type FileAuth struct {
	FileUUID userlib.UUID
	FileKey  []byte //accesses the file itself
	//this map is nil for top level authorized users. The owner will have access to a map of usernames to FileAuth UUIDs (for revokation)
	AuthUserMap map[string]TopAuthDetails

	IsOwner bool
}

//will be created whenever a file is created or shared to a new user by the user that the file has access granted to
//recipient will use data inside the invitation to create their own authpointer
//AuthPointers will be stored at uuid.FromBytes(userlib.Hash([]byte(userdata.Username + Filename))[:16])
//AuthPointers will be symmetrically encrypted using a user's own FileKey
type AuthPointer struct {
	AuthUUID userlib.UUID
	AuthKey  []byte
}

type FileChunk struct {
	NextChunkUUID userlib.UUID
	Content       []byte
	LastChunkUUID userlib.UUID
}

type Invitation struct {
	AuthUUID userlib.UUID
	AuthKey  []byte
}
type TopAuthDetails struct {
	AuthUUID userlib.UUID
	AuthKey  []byte
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {

	var userdata User
	userdataptr = &userdata
	userdata.Username = username
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username))[:16])
	if len(username) == 0 || err != nil {
		err = errors.New("Username cannot be an empty string!")
		return userdataptr, err //check all the errors later if they are correct
	}
	//check whether a user with the same username already exists
	_, usernametaken := userlib.DatastoreGet(userUUID)
	if usernametaken {
		err = errors.New("Username has already been taken!")
		return userdataptr, err
	}
	userdata.FileKey = userlib.RandomBytes(16)
	//create any other keys necessary to be stored in a user struct.
	//this creates the signing keys and stores them in keystore

	//FileAuth Public and Private Key

	//Digital Signature Public, private key pair
	signingkey, verifykey, err := userlib.DSKeyGen()

	userdata.SigningKey = signingkey

	err = userlib.KeystoreSet(username, verifykey)
	if err != nil {
		return userdataptr, err
	}

	invitationPubKey, invitationPrivKey, err := userlib.PKEKeyGen()
	userdata.InvitationKey = invitationPrivKey
	err = userlib.KeystoreSet(username+"pkey", invitationPubKey)
	if err != nil {
		return userdataptr, err
	}

	userBytes, err := json.Marshal(userdata)
	if err != nil {
		return userdataptr, err
	}
	//this creates the user encryption key
	marshalledPassword, err := json.Marshal(password)
	if err != nil {
		return userdataptr, err
	}
	marhsalledUsername, err := json.Marshal(username)
	if err != nil {
		return userdataptr, err
	}
	userEncryptionKey := userlib.Argon2Key(marshalledPassword, marhsalledUsername, 16)

	//Encrypt the userdata
	encryptedUser := userlib.SymEnc(userEncryptionKey, userlib.RandomBytes(16), userBytes)
	userMACKey, err := userlib.HashKDF(userEncryptionKey, []byte("mac"))
	userMACKey = userMACKey[:16]

	userMAC, err := userlib.HMACEval(userMACKey, encryptedUser)
	if err != nil {
		return userdataptr, err
	}
	userlib.DatastoreSet(userUUID, append(userMAC, encryptedUser...))
	return userdataptr, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {

	var userdata User
	userdataptr = &userdata
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])

	marshalledPassword, err := json.Marshal(password)
	if err != nil {
		return userdataptr, err
	}
	marhsalledUsername, err := json.Marshal(username)
	if err != nil {
		return userdataptr, err
	}
	userEncryptionKey := userlib.Argon2Key(marshalledPassword, marhsalledUsername, 16)
	userMACKey, err := userlib.HashKDF(userEncryptionKey, []byte("mac"))
	userMACKey = userMACKey[:16]
	if err != nil {
		return userdataptr, err
	}

	//get the user from the datastore
	encryptedUserWithMAC, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		err = errors.New("User does not exist!")
		return userdataptr, err
	}
	encryptedUser := encryptedUserWithMAC[64:]
	attachedUserMAC := encryptedUserWithMAC[:64]
	computedUserMAC, err := userlib.HMACEval(userMACKey, encryptedUser)
	if err != nil {
		return userdataptr, err

	}

	//verify the MAC

	userIntegrity := userlib.HMACEqual(attachedUserMAC, computedUserMAC)

	if !userIntegrity {
		err = errors.New("User has been tampered with!")
		return userdataptr, err
	}

	//decrypt the user

	marshalledUser := userlib.SymDec(userEncryptionKey, encryptedUser)
	err = json.Unmarshal(marshalledUser, &userdata)

	if err != nil {
		return userdataptr, err
	}
	//return the user data pointer

	return userdataptr, nil
}

//Files will be stored in chunks, using a randomly generated symmetric key and stored at a randomly generated UUID for the first chunk
//the first 16 bytes of the first chunk will contain the UUID of the last filechunk. the last 16 bytes of every file chunk will contain the UUID of the next chunk.
//the FileAuth object contains the UUID and key to the file.
func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	userAuthPointer, exists, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	//if user is overwriting an existing file
	if exists {
		userFileAuth, err := getFileAuth(userAuthPointer.AuthUUID, userAuthPointer.AuthKey)
		if err != nil {
			return err
		}
		fileUUID := userFileAuth.FileUUID

		fileKey := userFileAuth.FileKey

		chunkCount := (len(content) / 512) + 1

		curruuid := fileUUID
		initialUUID := fileUUID
		nextuuid := uuid.New()
		for i := 0; i < chunkCount; i++ {
			if len(content) < 512 {
				thisChunk := createFileChunk(nextuuid, content, uuid.New())
				err = storeFileChunk(curruuid, fileKey, thisChunk)
			} else {
				thisChunk := createFileChunk(nextuuid, content[:512], uuid.New())
				err = storeFileChunk(curruuid, fileKey, thisChunk)
				content = content[512:]
			}
			curruuid = nextuuid
			nextuuid = uuid.New()
			if err != nil {
				return err
			}
		}
		oldFirstFileChunk, _, err := getFileChunk(fileUUID, userFileAuth.FileKey)
		if err != nil {
			return err
		}
		newFirstFileChunk := createFileChunk(oldFirstFileChunk.NextChunkUUID, oldFirstFileChunk.Content, curruuid)

		storeFileChunk(initialUUID, userFileAuth.FileKey, newFirstFileChunk)

		return nil

	}

	//if user is creating a new file
	fileUUID := uuid.New()
	initialUUID := fileUUID
	fileKey := userlib.RandomBytes(16)
	if err != nil {
		return err
	}
	chunkCount := (len(content) / 512) + 1

	curruuid := fileUUID
	nextuuid := uuid.New()
	for i := 0; i < chunkCount; i++ {
		if len(content) < 512 {
			thisChunk := createFileChunk(nextuuid, content, uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
		} else {
			thisChunk := createFileChunk(nextuuid, content[:512], uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
			content = content[512:]
		}
		curruuid = nextuuid
		nextuuid = uuid.New()
		if err != nil {
			return err
		}
	}

	//creates owner's fileAuth using the file UUID and the file key
	ownerAuth := createFileAuth(initialUUID, fileKey, make(map[string]TopAuthDetails), true)
	ownerAuthUUID := uuid.New()
	//gets the user's public fileAuth encryption key stored in keystore
	//stores the user's fileAuth at a random UUID
	ownerFileAuthKey := userlib.RandomBytes(16)
	err = storeFileAuth(ownerAuthUUID, ownerFileAuthKey, ownerAuth)
	if err != nil {
		return err
	}

	oldFirstFileChunk, _, err := getFileChunk(fileUUID, ownerAuth.FileKey)
	newFirstFileChunk := createFileChunk(oldFirstFileChunk.NextChunkUUID, oldFirstFileChunk.Content, curruuid)

	storeFileChunk(initialUUID, ownerAuth.FileKey, newFirstFileChunk)

	//creates owner's authPointer
	//and stores the fileAuth's UUID and private fileAuth decryption key inside, then storing the authPointer object at a deterministic UUID
	//based on the username and filename in the owners namespace
	ownerAuthPointer := createAuthPointer(ownerAuthUUID, ownerFileAuthKey)

	storeAuthPointer(userdata.Username, filename, userdata.FileKey, ownerAuthPointer)

	return nil
}

func createAuthPointer(fileAuthUUID userlib.UUID, fileAuthKey []byte) (authPointer AuthPointer) {

	var thisAuthPointer AuthPointer
	thisAuthPointer.AuthUUID = fileAuthUUID
	thisAuthPointer.AuthKey = fileAuthKey

	return thisAuthPointer

}

//stores authPointer at a deterministic uuid determined by the username and filename, encrypted using symFileKey
func storeAuthPointer(username string, filename string, symFileKey []byte, authPointer AuthPointer) (err error) {
	authPointerUUID, err := uuid.FromBytes(userlib.Hash([]byte(username + filename + "ptr"))[:16])
	if err != nil {
		return err
	}

	marshalledPointer, err := json.Marshal(authPointer)
	if err != nil {
		return err
	}
	encryptedPointer := userlib.SymEnc(symFileKey, userlib.RandomBytes(16), marshalledPointer)
	if err != nil {
		return err
	}

	macKey, err := userlib.HashKDF(symFileKey, []byte("mac"))
	if err != nil {
		return err
	}
	macKey = macKey[:16]

	pointerMac, err := userlib.HMACEval(macKey, encryptedPointer)

	encryptedPointer = append(pointerMac, encryptedPointer...)

	userlib.DatastoreSet(authPointerUUID, encryptedPointer)
	return nil

}

//returns an authPointer object from dataStore when given the username, filename and symmetric file key of the user.
func getAuthPointer(username string, filename string, symFileKey []byte) (authPointer AuthPointer, exists bool, err error) {
	var authPointerData AuthPointer
	authPointerUUID, err := uuid.FromBytes(userlib.Hash([]byte(username + filename + "ptr"))[:16])
	if err != nil {
		return authPointerData, true, err
	}

	encryptedPointerWithMAC, exists := userlib.DatastoreGet(authPointerUUID)
	if !exists {
		//err = errors.New("AuthPointer does not exist!")
		return authPointerData, false, nil
	}
	encryptedAuthPointerMac := encryptedPointerWithMAC[:64]
	encryptedAuthPointer := encryptedPointerWithMAC[64:]
	//Check for Integrity
	macKey, err := userlib.HashKDF(symFileKey, []byte("mac"))
	if err != nil {
		return authPointerData, true, err
	}
	macKey = macKey[:16]

	computedPointerMac, err := userlib.HMACEval(macKey, encryptedAuthPointer)
	if err != nil {
		return authPointerData, true, err
	}
	pointerIntegrity := userlib.HMACEqual(encryptedAuthPointerMac, computedPointerMac)

	if !pointerIntegrity {
		err = errors.New("AuthPointer has been tampered with!")
		return authPointerData, true, err
	}

	marshalledPointer := userlib.SymDec(symFileKey, encryptedAuthPointer)

	err = json.Unmarshal(marshalledPointer, &authPointerData)
	if err != nil {
		return authPointerData, true, err
	}
	return authPointerData, true, nil

}

//if creating fileAuth for a non-owner, pass nil for userMap
//nextchunkUUID is the UUID of the next file chunk that has been stored at the last file chunk stored, put the first block of appended file contents there.
func createFileAuth(fileUUID userlib.UUID, fileKey []byte, userMap map[string]TopAuthDetails, isOwner bool) (fileAuth FileAuth) {
	var thisFileAuth FileAuth
	thisFileAuth.FileUUID = fileUUID
	thisFileAuth.FileKey = fileKey
	//thisFileAuth.NextChunkUUID = nextChunkUUID
	thisFileAuth.IsOwner = isOwner

	// AuthUserMap map[string][]byte
	thisFileAuth.AuthUserMap = userMap
	return thisFileAuth
}

//stores auth at uuid encrypted using pkey
func storeFileAuth(uuid userlib.UUID, pkey []byte, auth FileAuth) (err error) {

	marshalledAuth, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	encryptedAuth := userlib.SymEnc(pkey, userlib.RandomBytes(16), marshalledAuth)

	macKey, err := userlib.HashKDF(pkey, []byte("mac"))
	if err != nil {
		return err
	}
	macKey = macKey[:16]

	authMac, err := userlib.HMACEval(macKey, encryptedAuth)

	encryptedAuth = append(authMac, encryptedAuth...)

	userlib.DatastoreSet(uuid, encryptedAuth)
	return nil

}

//returns the FileAuth object from datastore when given the UUID and decryption key
func getFileAuth(uuid userlib.UUID, privkey []byte) (FileAuthObj FileAuth, err error) {
	var fileAuthData FileAuth
	encryptedFileAuthWithMac, exists := userlib.DatastoreGet(uuid)
	if !exists {
		err = errors.New("FileAuth does not exist!")
		return fileAuthData, err
	}

	encryptedAuthMac := encryptedFileAuthWithMac[:64]
	encryptedAuth := encryptedFileAuthWithMac[64:]
	//integrity MAC check
	macKey, err := userlib.HashKDF(privkey, []byte("mac"))
	if err != nil {
		return fileAuthData, err
	}
	macKey = macKey[:16]

	computedAuthMac, err := userlib.HMACEval(macKey, encryptedAuth)
	if err != nil {
		return fileAuthData, err
	}
	AuthIntegrity := userlib.HMACEqual(encryptedAuthMac, computedAuthMac)

	if !AuthIntegrity {
		err = errors.New("FileAuth has been tampered with!")
		return fileAuthData, err
	}

	decryptedFileAuth := userlib.SymDec(privkey, encryptedAuth)

	err = json.Unmarshal(decryptedFileAuth, &fileAuthData)
	if err != nil {
		return fileAuthData, err
	}
	return fileAuthData, nil
}

//func getFileAuth()

func createFileChunk(nextuuid userlib.UUID, content []byte, lastuuid userlib.UUID) (chunk FileChunk) {
	var thisChunk FileChunk
	thisChunk.NextChunkUUID = nextuuid
	thisChunk.Content = content
	thisChunk.LastChunkUUID = lastuuid
	return thisChunk
}

func storeFileChunk(uuid userlib.UUID, key []byte, chunk FileChunk) (err error) {
	marshalledChunk, err := json.Marshal(chunk)
	if err != nil {
		return err
	}
	encryptedChunk := userlib.SymEnc(key, userlib.RandomBytes(16), marshalledChunk)

	macKey, err := userlib.HashKDF(key, []byte("mac"))
	if err != nil {
		return err
	}
	macKey = macKey[:16]

	chunkMac, err := userlib.HMACEval(macKey, encryptedChunk)

	encryptedChunk = append(chunkMac, encryptedChunk...)
	userlib.DatastoreSet(uuid, encryptedChunk)
	return nil
}

func getFileChunk(uuid userlib.UUID, key []byte) (chunk FileChunk, exists bool, err error) {
	//if the file doesnt exist, do not error, return the fact the file doesnt exist.
	var chunkData FileChunk
	encryptedChunkWithMac, exists := userlib.DatastoreGet(uuid)
	if !exists {
		return chunkData, false, nil
	}

	encryptedChunkMac := encryptedChunkWithMac[:64]
	encryptedChunk := encryptedChunkWithMac[64:]
	//integrity MAC check
	macKey, err := userlib.HashKDF(key, []byte("mac"))
	if err != nil {
		return chunkData, true, err
	}
	macKey = macKey[:16]

	computedChunkMac, err := userlib.HMACEval(macKey, encryptedChunk)
	if err != nil {
		return chunkData, true, err
	}
	chunkIntegrity := userlib.HMACEqual(encryptedChunkMac, computedChunkMac)

	//Debug Messages
	// print("these are the macs \n")
	// print(encryptedChunkMac)
	// print(computedChunkMac)
	// print("\n chunk integrity ")
	// print(chunkIntegrity)

	if !chunkIntegrity {
		//print("IT HAS REACHED THE ERROR CONDITION")
		err = errors.New("file has been tampered with!")
		return chunkData, true, err
	}

	marshalledChunk := userlib.SymDec(key, encryptedChunk)

	err = json.Unmarshal(marshalledChunk, &chunkData)
	if err != nil {
		return chunkData, true, err
	}
	return chunkData, true, nil

}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	authPointer, _, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	if err != nil {
		return err
	}
	fileAuth, err := getFileAuth(authPointer.AuthUUID, authPointer.AuthKey)
	if err != nil {
		return err
	}

	firstChunk, exists, err := getFileChunk(fileAuth.FileUUID, fileAuth.FileKey)
	if exists == false {
		err = errors.New("File does not exist!")
		return err
	}
	if err != nil {
		return err
	}
	//appending
	appendStart := firstChunk.LastChunkUUID

	//Debug Statement
	//print("APPEND START")
	//print(firstChunk.LastChunkUUID.String())

	fileKey := fileAuth.FileKey

	chunkCount := (len(content) / 512) + 1

	curruuid := appendStart
	nextuuid := uuid.New()
	for i := 0; i < chunkCount; i++ {

		if len(content) < 512 {
			thisChunk := createFileChunk(nextuuid, content, uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
		} else {
			thisChunk := createFileChunk(nextuuid, content[:512], uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
			content = content[512:]
		}
		curruuid = nextuuid
		nextuuid = uuid.New()
		if err != nil {
			return err
		}
	}

	//creates a new fileAuth at the existing location with all existing values except update the last chunk UUID to the new one.
	//fileAuth = createFileAuth(fileAuth.FileUUID, fileAuth.FileKey, fileAuth.AuthUserMap, fileAuth.IsOwner)
	//storeFileAuth(authPointer.AuthUUID, authPointer.AuthKey, fileAuth)

	newFirstFileChunk := createFileChunk(firstChunk.NextChunkUUID, firstChunk.Content, curruuid)

	storeFileChunk(fileAuth.FileUUID, fileAuth.FileKey, newFirstFileChunk)

	return nil
}

//Error Cases: 1. filename does not exist in personal namespace of caller, 2. Integrity of the downloaded
//content cannot be verified, 3. Loading of file cannot succeed due to any other malicious action
//Need to cons
//Need to cons
//Need to cons
//Need to cons
func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	var filecontent []byte
	authPointer, _, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	if err != nil {
		return filecontent, err
	}
	fileAuth, err := getFileAuth(authPointer.AuthUUID, authPointer.AuthKey)
	if err != nil {
		return filecontent, err
	}
	//Get the file uuid from fileAuth
	//pass in file uuid to a getFileChunk call

	// type FileAuth struct {
	// 	FileUUID userlib.UUID
	// 	FileKey  []byte
	// 	//this map is nil for top level authorized users. The owner will have access to a map of usernames to FileAuth UUIDs (for revokation)
	// 	AuthUserMap   map[string][]byte
	// 	NextChunkUUID userlib.UUID
	// }

	//set chunk_var, exists to a function call getFileChunk with the current uuid
	chunk_var, exists, err := getFileChunk(fileAuth.FileUUID, fileAuth.FileKey)
	if err != nil {
		return filecontent, err
	}

	if exists == false {
		err = errors.New("File does not exist!")
		return filecontent, err
	}
	for exists != false {
		//	process the chunk_var
		filecontent = append(filecontent, chunk_var.Content...)
		chunk_var, exists, err = getFileChunk(chunk_var.NextChunkUUID, fileAuth.FileKey)
		if err != nil {
			return filecontent, err
		}
	}

	return filecontent, err
}

//for owners, create a new fileAuth, store the access details for the fileAuth inside a public key encrypted invitation
//for authorized users, store the access details for their own fileAuth inside a public key encrypted invitation
//the accepting user would create their own authPointer located in the corresponding UUID in their namespace containing these details.
func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	var invitePtr uuid.UUID
	authPointer, _, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	if err != nil {
		return invitePtr, err
	}
	fileAuth, err := getFileAuth(authPointer.AuthUUID, authPointer.AuthKey)
	if err != nil {
		return invitePtr, err
	}
	if fileAuth.IsOwner != true {
		//share location and key of existing user's fileAuth
		thisInvitation := createInvitationObject(authPointer.AuthUUID, authPointer.AuthKey)
		//encrypt the invite, store it in keystore at a random location, set that location to invitePtr, return.
		inviteUUID := uuid.New()
		publicKey, exists := userlib.KeystoreGet(recipientUsername + "pkey")
		if !exists {
			err = errors.New("Key does not exist!")
			return invitePtr, err
		}
		storeInvitationObject(inviteUUID, publicKey, thisInvitation, userdata.SigningKey)
		//invitePtr = inviteUUID
		return inviteUUID, nil
	}

	//create and store a new fileAuth here
	newFileAuth := createFileAuth(fileAuth.FileUUID, fileAuth.FileKey, make(map[string]TopAuthDetails), false)

	//encrypt this fileAuth using a new random key, add the key value pair of recipientUsername
	//and the location of the new fileAuth and this new random key inside the owners
	//original fileAuth.
	newFileAuthKey := userlib.RandomBytes(16)
	newFileAuthUUID := uuid.New()

	err = storeFileAuth(newFileAuthUUID, newFileAuthKey, newFileAuth)
	if err != nil {
		return invitePtr, err
	}

	inviteUUID := uuid.New()
	thisInvitation := createInvitationObject(newFileAuthUUID, newFileAuthKey)
	publicKey, exists := userlib.KeystoreGet(recipientUsername + "pkey")
	if !exists {
		err = errors.New("Key does not exist!")
		return invitePtr, err
	}
	storeInvitationObject(inviteUUID, publicKey, thisInvitation, userdata.SigningKey)

	//create a new hashmap containing all entries in the original hashmap + details of newFileAuth
	//then create a new fileAuth for the owner containing these details and store it.
	newTopAuthDetails := createTopAuthDetails(newFileAuthUUID, newFileAuthKey)
	newHashMap := fileAuth.AuthUserMap
	newHashMap[recipientUsername] = newTopAuthDetails

	newOwnerFileAuth := createFileAuth(fileAuth.FileUUID, fileAuth.FileKey, newHashMap, true)
	err = storeFileAuth(authPointer.AuthUUID, authPointer.AuthKey, newOwnerFileAuth)

	//invitePtr = inviteUUID
	return inviteUUID, err
}

func createTopAuthDetails(fileAuthUUID userlib.UUID, fileAuthKey []byte) (authDetails TopAuthDetails) {
	var thisTopAuthDetails TopAuthDetails
	thisTopAuthDetails.AuthUUID = fileAuthUUID
	thisTopAuthDetails.AuthKey = fileAuthKey

	return thisTopAuthDetails
}

func createInvitationObject(fileAuthUUID userlib.UUID, fileAuthKey []byte) (invitation Invitation) {

	var thisInvitation Invitation
	thisInvitation.AuthUUID = fileAuthUUID
	thisInvitation.AuthKey = fileAuthKey
	return thisInvitation
}
func storeInvitationObject(uuid userlib.UUID, publicKey userlib.PKEEncKey, invitation Invitation, signingKey userlib.DSSignKey) (err error) {
	marshalledInvitation, err := json.Marshal(invitation)
	if err != nil {
		return err
	}
	encryptedInvitation, err := userlib.PKEEnc(publicKey, marshalledInvitation)
	if err != nil {
		return err
	}

	signature, err := userlib.DSSign(signingKey, encryptedInvitation)
	if err != nil {
		return err
	}

	encryptedAndSignedInvitation := append(signature, encryptedInvitation...)
	userlib.DatastoreSet(uuid, encryptedAndSignedInvitation)

	return nil

}

func getInvitationObject(uuid userlib.UUID, dkey userlib.PKEDecKey, senderUsername string) (invitation Invitation, err error) {
	var invite Invitation
	encryptedInvitationWithSignature, exists := userlib.DatastoreGet(uuid)
	if !exists {
		err = errors.New("Invitation does not exist!")
		return invite, err
	}
	encryptedInvitation := encryptedInvitationWithSignature[256:]
	signature := encryptedInvitationWithSignature[:256]
	verifyKey, exists := userlib.KeystoreGet(senderUsername)
	if !exists {
		err = errors.New("Verification key does not exist!")
		return invite, err
	}
	err = userlib.DSVerify(verifyKey, encryptedInvitation, signature)
	if err != nil {
		return invite, err
	}

	marshalledInvitation, err := userlib.PKEDec(dkey, encryptedInvitation)

	err = json.Unmarshal(marshalledInvitation, &invite)
	if err != nil {
		return invite, err
	}

	return invite, nil

}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	//check if a file already exists in the recipients namespace
	authPointerUUID, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username + filename + "ptr"))[:16])
	if err != nil {
		return err
	}

	_, exists := userlib.DatastoreGet(authPointerUUID)
	if exists {
		err = errors.New("File Already Exists in Namespace!")
		return err
	}
	//create a new authPointer object for the user, using the data inside the invitation itself.
	invitation, err := getInvitationObject(invitationPtr, userdata.InvitationKey, senderUsername)
	if err != nil {
		return err
	}
	newAuthPointer := createAuthPointer(invitation.AuthUUID, invitation.AuthKey)
	err = storeAuthPointer(userdata.Username, filename, userdata.FileKey, newAuthPointer)
	if err != nil {
		return err
	}

	//debug statements
	//print("what is the FileAuthKey")
	//print(invitation.AuthKey)
	return nil

}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	//load the owner's fileAuth
	ownerAuthPointer, _, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	if err != nil {
		return err
	}
	ownerFileAuth, err := getFileAuth(ownerAuthPointer.AuthUUID, ownerAuthPointer.AuthKey)
	if err != nil {
		return err
	}
	//load the previous file
	userMap := ownerFileAuth.AuthUserMap

	recipientExists := false
	for names := range userMap {
		if names == recipientUsername {
			recipientExists = true
		}
	}

	if !recipientExists {
		err = errors.New("revoke recipient does not exist!")
		return err
	}

	fileContent, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}
	//file has been completely re encrypted, with access from all users revoked and a new authpointer and fileAuth created for the owner
	//with the authpointer stored at the original location since the filename and namespace are identical

	fileUUID := uuid.New()
	initialUUID := fileUUID
	fileKey := userlib.RandomBytes(16)
	if err != nil {
		return err
	}
	chunkCount := (len(fileContent) / 512) + 1

	curruuid := fileUUID
	nextuuid := uuid.New()
	for i := 0; i < chunkCount; i++ {
		if len(fileContent) < 512 {
			thisChunk := createFileChunk(nextuuid, fileContent, uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
		} else {
			thisChunk := createFileChunk(nextuuid, fileContent[:512], uuid.New())
			err = storeFileChunk(curruuid, fileKey, thisChunk)
			fileContent = fileContent[512:]
		}
		curruuid = nextuuid
		nextuuid = uuid.New()
		if err != nil {
			return err
		}
	}

	//creates owner's fileAuth using the file UUID and the file key
	ownerAuth := createFileAuth(initialUUID, fileKey, make(map[string]TopAuthDetails), true)
	ownerAuthUUID := uuid.New()
	//gets the user's public fileAuth encryption key stored in keystore
	//stores the user's fileAuth at a random UUID
	ownerFileAuthKey := userlib.RandomBytes(16)
	err = storeFileAuth(ownerAuthUUID, ownerFileAuthKey, ownerAuth)
	if err != nil {
		return err
	}

	oldFirstFileChunk, _, err := getFileChunk(fileUUID, ownerAuth.FileKey)
	newFirstFileChunk := createFileChunk(oldFirstFileChunk.NextChunkUUID, oldFirstFileChunk.Content, curruuid)

	storeFileChunk(initialUUID, ownerAuth.FileKey, newFirstFileChunk)

	//creates owner's authPointer
	//and stores the fileAuth's UUID and private fileAuth decryption key inside, then storing the authPointer object at a deterministic UUID
	//based on the username and filename in the owners namespace
	ownerTempAuthPointer := createAuthPointer(ownerAuthUUID, ownerFileAuthKey)

	storeAuthPointer(userdata.Username, filename, userdata.FileKey, ownerTempAuthPointer)

	//err = userdata.StoreFile(filename, fileContent)
	//if err != nil {
	//	return err
	//}

	//deletes the revoked top level user from userMap
	//scramble the original fileAuths
	for _, topAuthDetails := range userMap {
		updatedTopLevelFileAuth := createFileAuth(uuid.New(), userlib.RandomBytes(16), make(map[string]TopAuthDetails), false)
		err = storeFileAuth(topAuthDetails.AuthUUID, topAuthDetails.AuthKey, updatedTopLevelFileAuth)
		if err != nil {
			return err
		}
	}

	//Debug Statements
	// print("before delete")
	// for key, _ := range userMap {
	// 	print(key)
	// 	print("\n")
	// }

	delete(userMap, recipientUsername)

	// print("after delete")
	// for key, _ := range userMap {
	// 	print(key)
	// 	print("\n")
	// }

	//overwrite the old file chunk with a new garbled filechunk, encrypted using a throwaway key
	//randomChunk := createFileChunk(uuid.New(), userlib.RandomBytes(5), uuid.New())
	//err = storeFileChunk(ownerFileAuth.FileUUID, userlib.RandomBytes(16), randomChunk)

	userlib.DatastoreDelete(ownerFileAuth.FileUUID)

	ownerNewAuthPointer, _, err := getAuthPointer(userdata.Username, filename, userdata.FileKey)
	if err != nil {
		return err
	}
	ownerNewFileAuth, err := getFileAuth(ownerNewAuthPointer.AuthUUID, ownerNewAuthPointer.AuthKey)
	if err != nil {
		return err
	}

	for _, topAuthDetails := range userMap {
		updatedTopLevelFileAuth := createFileAuth(ownerNewFileAuth.FileUUID, ownerNewFileAuth.FileKey, make(map[string]TopAuthDetails), false)
		err = storeFileAuth(topAuthDetails.AuthUUID, topAuthDetails.AuthKey, updatedTopLevelFileAuth)
		if err != nil {
			return err
		}
	}

	return nil
}
