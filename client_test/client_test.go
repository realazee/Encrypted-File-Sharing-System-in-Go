package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	strconv "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const contentFour = "In the past years, the staple hip-hop number Gucci Gang has become a fixture of popular culture, its influence extending past the maturing Generation Z to all aspects of daily life and conduct. Like the era of Reaganism proceeding it, despite political inclinations and decrying of “moral decay” by the white ruling class, “Gucci Gang” appears to have become a counter-culture display of multiculturalism and individualism in an isolationist world. The author of this work, Lil Pump, has put front and center in his bold oeuvre the psychological strain and nebulousness of individuality in an age where tribalism and cultural identity of individual cliques is encouraged. He furthermore elaborates that his strain has degraded the meaning and nuance in his life. 	The corrupted state of the young generation is first and foremost, implied through the song’s seemingly paradoxical lyrics. The repetition of the phrase “Gucci Gang” not only re-addresses the title but also contains seemingly contradictory life values and imagery. The brand Gucci is associated with the wealthier classes of society, as only this privileged group has the ability to purchase such luxury esthetics with the well-known brand Gucci. On the other hand, the word “gang” evocates the stereotype of an impoverished, vandalized, and lawless area filled with danger and insecurity - a message that directly goes against the luxurious quality of Gucci and the top one percent. This diametric, split messaging shows how modern American society has influenced its people so that they have a similar split identity.  What’s shown above is symbolized by how, despite circumstances, we all have an insatiable drive for comfort and success, embodied by “Gucci” and the promise and allure it holds. On the opposite side of the spectrum, the word “gang” nods to the street culture which is inspired by life inside a gang; for example rapping, “sagging”, and unfortunately, drugs; which shows the fashion trends that are considered “hood.”  Thus, the author’s thoughtful placement of symbolism tells people of the mixed identity that modern teens possess, such as the contrast between the luxury brand Gucci and “the hood” which is traditionally associated with poorer communities.  	After the first line of the song that introduces people to his ideas, the singer then moves on to the darker side of superficiality, where he begins to elaborate on the psychological toll of teenage indecisiveness. He illustrates to the listeners of the song how he routinely “forgets the name” of those he interacts with on a daily basis, implying that while his materialistic wishes give him solace, they have also made him blind to the beauties and intricacy of the world around him. This is further reinforced by his staunch belief that he “can’t buy no wedding ring,” as he would “rather go and buy Balmains.” Despite the author’s acknowledgment that love and beauty do exist around him, he rejects their presence so in order to patronize “Gucci” and keep his own darker side at bay. This emptiness ultimately consumes one’s character over the duration of this song, leading him to denigrate everything from “airlines” to “cigarettes,” as these can no longer satisfy his intrinsic need for decadence. 	Throughout the song, the singer repeats the names of drugs multiple times, and in the way in which they are described such as: “love[to] do cocaine,” “breath smell like some cigarettes,” “still sell that meth.” The repetition of the influence of drugs is a way by the singer to express the concern for the serious addiction problems teenagers in the modern era face. No matter if it’s social media, games, or literal drugs, more and more teenagers are now addicted to the innumerable amounts of stimulus that they take in every day; just like literal drugs that harm one's body, addictions to other less material content is still incredibly harmful. These expressions also bring in the factor of pressure from their surroundings, as the subjects who are associated with drugs in the lyrics are all close to the singer as seen in “My [girl] love do cocaine”, “my grandma take[s] meds”. The singer infers that in the epidemic of teenage addiction, the teens are not the only ones to blame. Their peers and even the adults that are supposed to act as role-models for said teens have all fallen under the spell of addiction with no other examples to refer to, they naturally fall into the same path along with the addicts.  	Finally, the rhythm and the tone of the song, along with the music that ultimately makes it a song, but not a poem, also spread a message; a message that is not intended by the singer. The swear words that are almost in every single line of the lyrics, the loud music accompanying the lyrics (that, to be honest, basically frowns out any audible words sung by the singer), and the close rhymes that no matter how one puts it, do not work, all show the desperate state that the teenage art scene is in. In the twenty-first century, music is no longer something that is used to relax and calm the mind, it is no longer supposed to be compositionally interesting, figurative, and a carefully crafted work of art that the artists pour their heart and soul into. It is not the grand orchestral classical music of old; no longer the national anthems that unite entire nations; no longer the patriotic tones of the 1910’s; no longer the versatile, original, and symbolic music of David Bowie. Now songs like Gucci Gang, where lines do not rhyme, where words are meaningless; where not only is the singer described as doing drugs, but it’s “Me and my grandma take meds”; where love is reduced to Gucci, Balmain, and cocaine, has taken the main stage.  	Thus, throughout the course of Gucci Gang, we see the narrator transition from being torn between his competing desires to his eventual hatred of everything not fitting into his vain appeals, in the resolution. In a striking, eerily thoughtless display, Lil Pump thus illustrates the emptiness of both materialism and lust with his experiences in a contemporary, relevant fashion. Gucci Gang thus calls upon all Americans to rise up above their wishful, selfish, and superficial wants and realize that happiness is not tribal, nor is it based on empty happiness or bloodshed. But instead that the world is a place where people don’t do drugs for a living, where women aren't female dogs, where buying Gucci and Balmain don’t represent love, and where diverse populations do come together and care about each other."

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	dorisFile := "dorisFile.txt"
	eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			//**********Added Test**********
			userlib.DebugMsg("Checking that bobDesktop sees expected file data.")
			bobdata, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobdata).To(Equal([]byte(contentOne)))
			//**********Added Test**********

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			//self added test
			//userlib.DebugMsg("Checking that Bob/Charles do not have access to the contents of the file")
			//bobFile, err := bob.LoadFile(bobFile)
			//Expect(err).ToNot(BeNil())
			//self added test

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		//[Design Req 3.6.4] The client must enfore that there is only a single copy of a file.
		Specify("Design Req 3.6.4: The client must enfore that there is only a single copy of a file.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			//**********Added Test**********
			userlib.DebugMsg("Checking that bob sees expected file data.")
			bobdata, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobdata).To(Equal([]byte(contentOne)))
			//**********Added Test**********

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alice sees expected, modified file data.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)))

		})

		//[Design Req 3.5.5] The client may leak any information except filenames, lengths of filenames, file contents,
		// and file sharing invitations
		Specify("Design Req 3.5.5: The client may leak any information except filenames, lengths of filenames, file contents, and file sharing invitations.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice storing file %s", aliceFile)
			err = alice.StoreFile(aliceFile, userlib.RandomBytes(10000000))
			Expect(err).To(BeNil())

			userlib.DebugMsg("charles storing file %s", charlesFile)
			err = charles.StoreFile(charlesFile, userlib.RandomBytes(10000000))
			Expect(err).To(BeNil())

			datastore := userlib.DatastoreGetMap()

			//tamper with file contents
			for key, value := range datastore {

				if len(value) > 512 {
					print("File has been tampered in datastore")
					datastore[key] = userlib.RandomBytes(len(value))
				}
				//print("\n")
				//print("key from datastore")
				//print(key.String())
				//print(" ")
				//print("value from datastore")
				//print(value)
				//print("\n")

			}
			userlib.DebugMsg("alice loading tampered file")
			_, err = alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("charles appending to tampered file")
			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

		})

		//[Design Req 3.5.2] The client must ensure the integrity of filenames

	})

	//more revoke functionality test, checking that other users the file has been shared with other than users
	//that have been revoked still have access to the file
	Specify("Testing Advanced Revoke Functionality", func() {
		userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		charles, err = client.InitUser("charles", defaultPassword)
		Expect(err).To(BeNil())

		doris, err = client.InitUser("doris", defaultPassword)
		Expect(err).To(BeNil())

		eve, err = client.InitUser("eve", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentFour)
		alice.StoreFile(aliceFile, []byte(contentFour))

		userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice creating invite for Doris for file %s, and Doris accepting invite under name %s.", aliceFile, dorisFile)

		invite, err = alice.CreateInvitation(aliceFile, "doris")
		Expect(err).To(BeNil())

		err = doris.AcceptInvitation("alice", invite, dorisFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Alice can still load the file.")
		data, err := alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Checking that Bob can load the file.")
		data, err = bob.LoadFile(bobFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Checking that Doris can load the file.")
		data, err = doris.LoadFile(dorisFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
		invite, err = bob.CreateInvitation(bobFile, "charles")
		Expect(err).To(BeNil())

		err = charles.AcceptInvitation("bob", invite, charlesFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Charles can load the file.")
		data, err = charles.LoadFile(charlesFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Doris creating invite for Eve for file %s, and Eve accepting invite under name %s.", dorisFile, eveFile)
		invite, err = doris.CreateInvitation(dorisFile, "eve")
		Expect(err).To(BeNil())

		err = eve.AcceptInvitation("doris", invite, eveFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Eve can load the file.")
		data, err = eve.LoadFile(eveFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
		err = alice.RevokeAccess(aliceFile, "bob")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Alice can still load the file.")
		data, err = alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Checking that Doris can still load the file.")
		data, err = doris.LoadFile(dorisFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		userlib.DebugMsg("Checking that Eve can still load the file.")
		data, err = eve.LoadFile(eveFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

		//self added test
		//userlib.DebugMsg("Checking that Bob/Charles do not have access to the contents of the file")
		//bobFile, err := bob.LoadFile(bobFile)
		//Expect(err).ToNot(BeNil())
		//self added test

		userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
		_, err = bob.LoadFile(bobFile)
		Expect(err).ToNot(BeNil())

		_, err = charles.LoadFile(charlesFile)
		Expect(err).ToNot(BeNil())

		userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
		err = bob.AppendToFile(bobFile, []byte(contentTwo))
		Expect(err).ToNot(BeNil())

		err = charles.AppendToFile(charlesFile, []byte(contentTwo))
		Expect(err).ToNot(BeNil())
	})

	//helper to measure bandwidth
	measureBandwidth := func(probe func()) (bandwidth int) {
		before := userlib.DatastoreGetBandwidth()
		probe()
		after := userlib.DatastoreGetBandwidth()
		return after - before
	}

	//testing storing and appending long files and corresponding bandwidth
	Specify("Testing storing and appending very long files and corresponding bandwidth", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentFour)
		err = alice.StoreFile(aliceFile, []byte(contentFour))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Appending file data: %s", contentTwo)
		bwFirstAppend := measureBandwidth(func() {
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
		})
		Expect(err).To(BeNil())
		Expect(bwFirstAppend).To(BeNumerically("<", 3000))

		userlib.DebugMsg("Appending file data: %s", contentFour)
		err = alice.AppendToFile(aliceFile, []byte(contentFour))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Appending file data: %s", contentThree)

		bwSecondAppend := measureBandwidth(func() {
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
		})
		Expect(err).To(BeNil())
		Expect(bwSecondAppend).To(BeNumerically("<", 3000))

		userlib.DebugMsg("Loading file...")
		data, err := alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour + contentTwo + contentFour + contentThree)))

		userlib.DebugMsg("Appending Extremely Long Random File Data")
		err = alice.AppendToFile(aliceFile, userlib.RandomBytes(10000000))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Appending file data: %s", contentOne)

		bwThirdAppend := measureBandwidth(func() {
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
		})
		Expect(err).To(BeNil())
		Expect(bwThirdAppend).To(BeNumerically("<", 3000))

		//perform 10000 appends
		i := 0
		for i < 10000 {
			err = alice.AppendToFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			i++
		}

		bwFourthAppend := measureBandwidth(func() {
			err = alice.AppendToFile(aliceFile, []byte(contentOne))
		})
		Expect(err).To(BeNil())
		Expect(bwFourthAppend).To(BeNumerically("<", 3000))

	})

	//Test that inviting a user who does not exist throws an error.
	Specify("Test that inviting a user who does not exist throws an error .", func() {
		userlib.DebugMsg("Initializing user Alice")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Frank, who does not exist.")
		_, err := alice.CreateInvitation(aliceFile, "frank")
		Expect(err).ToNot(BeNil())

	})
	// Specify("Test that inviting an already authorized user throws an error.", func() {
	// 	userlib.DebugMsg("Initializing users Alice and Charles")
	// 	alice, err = client.InitUser("alice", defaultPassword)
	// 	Expect(err).To(BeNil())

	// 	charles, err = client.InitUser("charles", defaultPassword)
	// 	Expect(err).To(BeNil())

	// 	userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentThree)
	// 	err = alice.StoreFile(aliceFile, []byte(contentThree))
	// 	Expect(err).To(BeNil())

	// 	userlib.DebugMsg("alice creating invite for Charles")
	// 	_, err := alice.CreateInvitation(aliceFile, "charles")
	// 	Expect(err).ToNot(BeNil())

	// })
	Specify("Test that client supports a user with password equal to an empty string.", func() {
		userlib.DebugMsg("Initializing user Doris")
		doris, err = client.InitUser("doris", "")
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris storing file %s with content: %s", dorisFile, contentFour)
		err = doris.StoreFile(dorisFile, []byte(contentFour))
		Expect(err).To(BeNil())
	})

	//Test that revoking a user who does not exist throws and error.
	//Test that revoking a user .
	Specify("Test that revoking a user from a nonexistent file errors.", func() {
		userlib.DebugMsg("Initializing users Alice")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice revoking Bob's access from %s, a nonexistent file", charlesFile)
		err = alice.RevokeAccess(charlesFile, "bob")
		Expect(err).ToNot(BeNil())
	})

	//check that different users can store files using the same name

	Specify("Testing same file name under different users", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Initializing user Bob.")
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = bob.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice Loading file...")
		data, err := alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentOne)))

		userlib.DebugMsg("Bob Loading file...")
		data, err = bob.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentOne)))

	})
	//check that file names can be empty strings
	Specify("Testing empty string file name", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile("", []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice Loading file...")
		data, err := alice.LoadFile("")
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentOne)))

	})

	//check that the number of entries in keystore does not increase based on number of files.
	Specify("Testing keystore key count does not increase with number of files", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		i := 0
		for i < 1000 {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(strconv.Itoa(i), []byte(contentOne))
			Expect(err).To(BeNil())
			i++
		}

		userlib.KeystoreGetMap()
		Expect(len(userlib.KeystoreGetMap())).To(BeNumerically("<", 10))

	})

	//testing using incorrect password
	Specify("Testing using an incorrect password.", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Getting user Alice.")
		aliceLaptop, err = client.GetUser("alice", "Not the defaultPassword")
		Expect(err).ToNot(BeNil())

	})
	//testing logging into user that doesnt exist.
	Specify("Testing logging into a user that doesn't exist", func() {

		userlib.DebugMsg("Getting user Alice.")
		aliceLaptop, err = client.GetUser("alice", defaultPassword)
		Expect(err).ToNot(BeNil())

	})
	//Testing that Loading a File that does not exist errors
	Specify("Testing that Loading a File that does not exist errors", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice Loading file...")
		_, err := alice.LoadFile("fakeFile")
		Expect(err).ToNot(BeNil())

	})
	//Testing that a user accepting an invitation for a file that already exists in their personal namespace
	//errors
	Specify("Testing that a user accepting an invitation for a file that already exists in their personal namespace errors", func() {
		userlib.DebugMsg("Initializing user Doris")
		doris, err = client.InitUser("doris", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris storing file %s with content: %s", dorisFile, contentOne)
		err = doris.StoreFile(dorisFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris creating invite for bob to fakefile, which does not exist")
		invite, err := doris.CreateInvitation(dorisFile, "bob")
		Expect(err).To(BeNil())

		err = bob.AcceptInvitation("doris", invite, bobFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris storing file %s with content: %s", aliceFile, contentTwo)
		err = doris.StoreFile(aliceFile, []byte(contentTwo))
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris creating invite for bob to fakefile, which does not exist")
		invite2, err := doris.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		err = bob.AcceptInvitation("doris", invite2, bobFile)
		Expect(err).ToNot(BeNil())
	})

	//Testing that Creating an Invite to a File that does not exist errors
	Specify("Testing that Creating an Invite to a File that does not exist errors.", func() {
		userlib.DebugMsg("Initializing user Doris")
		doris, err = client.InitUser("doris", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris storing file %s with content: %s", dorisFile, contentOne)
		err = doris.StoreFile(dorisFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("doris creating invite for bob to fakefile, which does not exist")
		_, err := alice.CreateInvitation("fakefile", "bob")
		Expect(err).ToNot(BeNil())

	})

	//Testing that Appending to a File that does not exist errors
	Specify("Testing that Appending to a File that does not exist errors", func() {
		userlib.DebugMsg("Initializing user Charles.")
		charles, err = client.InitUser("charles", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = charles.StoreFile(charlesFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Charles Loading file...")
		_, err := charles.LoadFile(charlesFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Charles Appending file...")
		err = charles.AppendToFile("fakeFile", []byte(contentTwo))
		Expect(err).ToNot(BeNil())

	})

	//testing shared user file overwriting
	Specify("Testing shared user file overwriting.", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Bob.")
		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Bob can load the file.")
		data, err := bob.LoadFile(bobFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentOne)))

		userlib.DebugMsg("checking that bob can overwrite the file")
		err = bob.StoreFile(bobFile, []byte(contentFour))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that alice can load the file.")
		data, err = alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentFour)))

	})

	//testing accessing a file that a user doesnt have access to.
	Specify("Testing accessing file without being invited", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Bob cannot load the file.")
		data, err := bob.LoadFile(aliceFile)
		Expect(err).ToNot(BeNil())
		Expect(data).ToNot(Equal([]byte(contentOne)))

	})

	Specify("Testing revoking a file after an invitation has been sent but not accepted", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Bob.")
		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Alice revoking Bob's access from %s", aliceFile)
		err = alice.RevokeAccess(aliceFile, "bob")
		Expect(err).To(BeNil())

		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Bob cannot load the file.")
		data, err := bob.LoadFile(aliceFile)
		Expect(err).ToNot(BeNil())
		Expect(data).ToNot(Equal([]byte(contentOne)))
	})

	Specify("Testing creating a user with an empty username", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("", defaultPassword)
		Expect(err).ToNot(BeNil())

	})
	Specify("Testing creating a user that already exists", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Initializing user bob.")
		bob, err = client.InitUser("alice", defaultPassword)
		Expect(err).ToNot(BeNil())

	})

	Specify("If the caller to AcceptInvitation is unable to verify that the secure file Invitation pointed to by the given invitationPtr was created by senderUsername, throw an error.", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Initializing user bob.")
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentFour)
		err = alice.StoreFile(aliceFile, []byte(contentFour))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Bob.")
		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		err = bob.AcceptInvitation("charles", invite, bobFile)
		Expect(err).ToNot(BeNil())
	})

	Specify("Test for revoking a user from a file that has not been shared with the user", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		doris, err = client.InitUser("doris", "_weoij23")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentTwo)
		err = doris.StoreFile(dorisFile, []byte(contentTwo))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Doris revoking Alice's access from %s, which has not been shared with Alice", dorisFile)
		err = doris.RevokeAccess(dorisFile, "alice")
		Expect(err).ToNot(BeNil())
	})

	Specify("Test that overwriting a file does not change who it is shared with.", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Bob.")
		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentTwo)
		err = alice.StoreFile(aliceFile, []byte(contentTwo))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Bob can load the file.")
		data, err := bob.LoadFile(bobFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentTwo)))

	})
	Specify("Test that tampering with AuthPointers causes an error", func() {
		userlib.DebugMsg("Initializing user Alice.")
		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		//tampering with the authPointer

		m := userlib.DatastoreGetMap()

		pointerUUID, err := uuid.FromBytes(userlib.Hash([]byte("alice" + aliceFile + "ptr"))[:16])
		Expect(err).To(BeNil())

		m[pointerUUID] = userlib.RandomBytes(256)
		userlib.DebugMsg("Checking that Alice loading the file will error.")
		data, err := alice.LoadFile(aliceFile)
		Expect(err).ToNot(BeNil())
		Expect(data).ToNot(Equal([]byte(contentOne)))

	})
	Specify("Test that tampering a User is detected.", func() {
		userlib.DebugMsg("Initializing user Charles.")
		charles, err = client.InitUser("charles", defaultPassword)
		Expect(err).To(BeNil())

		m := userlib.DatastoreGetMap()

		userUUID, err := uuid.FromBytes(userlib.Hash([]byte("charles"))[:16])

		m[userUUID] = userlib.RandomBytes(256)
		userlib.DebugMsg("Checking that getting tampered user errors.")
		_, err = client.GetUser("charles", defaultPassword)
		Expect(err).ToNot(BeNil())

	})
	Specify("Test that creating a fake invitation is detected", func() {
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Mallory creating fake invite for Bob.")
		m := userlib.DatastoreGetMap()

		invite := uuid.New()
		m[invite] = userlib.RandomBytes(256)

		userlib.DebugMsg("Bob accepting fake invite from Alice under filename %s.", bobFile)
		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).ToNot(BeNil())

	})

	Specify("Test that modifying an invitation is detected", func() {
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("alice creating invite for Bob.")
		invite, err := alice.CreateInvitation(aliceFile, "bob")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Mallory tampering with invite for Bob.")
		m := userlib.DatastoreGetMap()

		m[invite] = append(m[invite], userlib.RandomBytes(16)...)

		userlib.DebugMsg("Bob accepting altered invite from Alice under filename %s.", bobFile)
		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).ToNot(BeNil())

	})

	Specify("Test that usernames are case-sensitive", func() {
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())
		bob, err = client.InitUser("Alice", "differentpass")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Getting user Alice.")
		aliceLaptop, err = client.GetUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Getting user Bob.")
		bob, err = client.GetUser("Alice", "differentpass")
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentTwo)
		err = bob.StoreFile(aliceFile, []byte(contentTwo))
		Expect(err).To(BeNil())

		userlib.DebugMsg("Checking that Alice can load the file.")
		data, err := alice.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentOne)))

		userlib.DebugMsg("Checking that Bob can load the file.")
		data, err = bob.LoadFile(aliceFile)
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte(contentTwo)))

	})
	Specify("Test that getting a user that does not exist errors", func() {
		userlib.DebugMsg("Getting user Charles.")
		charles, err = client.GetUser("charles", defaultPassword)
		Expect(err).ToNot(BeNil())
	})

	Specify("Fuzz Datastore Tampering", func() {

		i := 0
		for i < 20 {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			tampered := false

			m := userlib.DatastoreGetMap()
			for key := range m {
				if GinkgoRandomSeed() > 0 {
					tampered = true
					m[key] = userlib.RandomBytes(256)
				}
			}
			data, err := alice.LoadFile(aliceFile)
			if tampered {
				print("TAMPERED")
				Expect(err).ToNot(BeNil())
				Expect(data).ToNot(Equal([]byte(contentOne)))
			} else {
				print("NOT MESSED WITH")
				Expect(err).To(BeNil())
				Expect(data).To(Equal([]byte(contentOne)))
			}
			userlib.DatastoreClear()
			userlib.KeystoreClear()
			i++

		}

	})
	Specify("Testing accepting an invitation that doesn't exist", func() {
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Storing file data: %s", contentOne)
		err = alice.StoreFile(aliceFile, []byte(contentOne))
		Expect(err).To(BeNil())

		invite := uuid.New()
		err = bob.AcceptInvitation("alice", invite, bobFile)
		Expect(err).ToNot(BeNil())

	})
	Specify("Testing swapping users in datastore", func() {
		userlib.DebugMsg("Initializing user Alice.")
		alice, err = client.InitUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Initializing user Bob.")
		bob, err = client.InitUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Getting user Alice.")
		alice, err = client.GetUser("alice", defaultPassword)
		Expect(err).To(BeNil())

		userlib.DebugMsg("Getting user Bob.")
		bob, err = client.GetUser("bob", defaultPassword)
		Expect(err).To(BeNil())

		m := userlib.DatastoreGetMap()

		aliceLocation, err := uuid.FromBytes(userlib.Hash([]byte("alice"))[:16])
		Expect(err).To(BeNil())
		bobLocation, err := uuid.FromBytes(userlib.Hash([]byte("bob"))[:16])
		Expect(err).To(BeNil())

		aliceUser := m[aliceLocation]
		bobUser := m[bobLocation]

		m[aliceLocation] = bobUser
		m[bobLocation] = aliceUser

		userlib.DebugMsg("Getting user Alice.")
		alice, err = client.GetUser("alice", defaultPassword)
		Expect(err).ToNot(BeNil())

		userlib.DebugMsg("Getting user Bob.")
		bob, err = client.GetUser("bob", defaultPassword)
		Expect(err).ToNot(BeNil())

	})

})

//multi-session testing
