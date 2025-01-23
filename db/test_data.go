package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

const (
	johnCogId     = "00001"
	kevinCogId    = "00002"
	blakeCogId    = "00003"
	martyCogId    = "00004"
	assetsDirName = "/home/john-marren/Projects/deepfried/static/assets/"
)

func InitTestData(ctx context.Context) error {
	err := initTestUsers(ctx)
	if err != nil {
		return fmt.Errorf("error initTestUsers: %s", err)
	}
	err = initTestAudioFiles(ctx)

	if err != nil {
		return fmt.Errorf("error initTestAudioFiles: %s", err)
	} else {
		fmt.Println("initTestAudioFiles finished without error")
	}

	return nil
}

type insertTestUserData struct {
	CognitoID string
	Username  string
	Bio       pgtype.Text
}

func initTestUsers(ctx context.Context) error {
	testUsers := []insertTestUserData{
		{
			CognitoID: johnCogId,
			Username:  "john",
			Bio: pgtype.Text{
				String: "just a man with a plan",
				Valid:  true,
			},
		},
		{
			CognitoID: kevinCogId,
			Username:  "wonderlust",
			Bio: pgtype.Text{
				String: "Imma bad bitch",
				Valid:  true,
			},
		},
		{
			CognitoID: blakeCogId,
			Username:  "blakefoster",
			Bio: pgtype.Text{
				String: "",
				Valid:  false,
			},
		},
		{
			CognitoID: martyCogId,
			Username:  "lovejames",
			Bio: pgtype.Text{
				String: "ankles deep in fun, always",
				Valid:  true,
			},
		},
	}

	tx, err := Dbtx.Begin(ctx)
	util.EMsg(err, "beginning db tx")
	qtx := Query.WithTx(tx)

	for _, user := range testUsers {

		userId, err := qtx.CreateUserTest(ctx, sqlc.CreateUserTestParams{
			CognitoID: user.CognitoID,
			Username:  user.Username,
		})
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error inserting test users: %s", err)
		}
		err = qtx.CreateProfileTest(ctx, sqlc.CreateProfileTestParams{
			UserID: userId,
			Bio:    user.Bio,
		})
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error inserting test profile: %s", err)
		}
	}

	tx.Commit(ctx)
	return nil
}

func copyFileLocations(userId pgtype.UUID, username string) {

	oldProfilePhoto := fmt.Sprintf("%sprofile_photos/%s", assetsDirName, username)
	newProfilePhoto := fmt.Sprintf("%sprofile_photos/%s", assetsDirName, util.UuidString(userId))

	_, err := os.Stat(newProfilePhoto)
	if err != nil {
		cmd := exec.Command("cp", oldProfilePhoto, newProfilePhoto)
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%s %s", oldProfilePhoto, newProfilePhoto))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Printf(string(output))
	}

	oldArtwork := fmt.Sprintf("%sartwork/%s", assetsDirName, username)
	newArtwork := fmt.Sprintf("%sartwork/%s/", assetsDirName, util.UuidString(userId))

	_, err = os.Stat(newArtwork)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		cmd := exec.Command("cp", "-R", oldArtwork, newArtwork)
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%s %s", oldArtwork, newArtwork))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf(string(output))
	}

}

func initTestAudioFiles(ctx context.Context) error {
	john, err := Query.GetUserWithCognitoId(ctx, "00001")
	if err != nil {
		return fmt.Errorf("error getting john: %s\n", err)
	}
	kevin, err := Query.GetUserWithCognitoId(ctx, kevinCogId)
	if err != nil {
		return fmt.Errorf("error getting kevin: %s", err)
	}
	blake, err := Query.GetUserWithCognitoId(ctx, "00003")
	if err != nil {
		return fmt.Errorf("error getting blake: %s\n", err)
	}
	marty, err := Query.GetUserWithCognitoId(ctx, martyCogId)
	if err != nil {
		return fmt.Errorf("error getting marty: %s", err)
	}

	copyFileLocations(john.ID, john.Username)
	copyFileLocations(kevin.ID, kevin.Username)
	copyFileLocations(blake.ID, blake.Username)
	copyFileLocations(marty.ID, marty.Username)

	// newLocation := fmt.Sprintf("/home/john-marren/Projects/loop-project/static/assets/profile_photos/%s", UuidString(kevin.ID))
	// prevLocation := "/home/john-marren/Projects/loop-project/static/assets/profile_photos/wonderlust"

	testAudioFiles := []sqlc.AddTestAudioFileParams{
		{
			UserID:              kevin.ID,
			Title:               "BALLROOM",
			Bpm:                 67,
			PlaybackSeconds:     115,
			FileSize:            3300,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "BALLROOM",
			AudioSrc:            "BALLROOM",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "MAJDAG",
			Bpm:                 140,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "MAJDAG",
			AudioSrc:            "MAJDAG",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "McLovin",
			Bpm:                 160,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "McLovin",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Nightmaze",
			Bpm:                 160,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "Nightmaze.png",
			AudioSrc:            "Nightmaze",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "SEX",
			Bpm:                 108,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "SEX",
			AudioSrc:            "SEX",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Heaven",
			Bpm:                 108,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "Heaven.png",
			AudioSrc:            "Heaven",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Different",
			Bpm:                 108,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "Different",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "TV Girl",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "TV Girl",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Ripe",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "Ripe",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "6 Feet Deep",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "6 Feet Deep",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Jubilee",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "McLovin",
			AudioSrc:            "Jubilee",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Fly With Me",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "fly_with_me.jpg",
			AudioSrc:            "Fly With Me",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Horizon",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "horizon.jpg",
			AudioSrc:            "Horizon",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Beverly",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "beverly.jpeg",
			AudioSrc:            "Beverly",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              kevin.ID,
			Title:               "Modest",
			Bpm:                 96,
			PlaybackSeconds:     69,
			FileSize:            2000,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "Modest.jpg",
			AudioSrc:            "Modest",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              marty.ID,
			Title:               "Swifty",
			Bpm:                 145,
			PlaybackSeconds:     60,
			FileSize:            10600,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "Swifty",
			AudioSrc:            "Swifty",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              john.ID,
			Title:               "Blueberry Fields",
			Bpm:                 145,
			PlaybackSeconds:     60,
			FileSize:            10600,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "blueberry_fields_artwork",
			AudioSrc:            "Blueberry Fields",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              blake.ID,
			Title:               "Bad Dogs",
			Bpm:                 145,
			PlaybackSeconds:     60,
			FileSize:            10600,
			VisArr:              []int32{1, 2, 3, 4, 44},
			ArtworkSrc:          "blakefoster",
			AudioSrc:            "Bad Dogs",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
		{
			UserID:              blake.ID,
			Title:               "Star Projector",
			Bpm:                 145,
			PlaybackSeconds:     60,
			FileSize:            10600,
			VisArr:              []int32{1, 2, 3, 4, 44},
			AudioSrc:            "Star Projector",
			MusicalKey:          "C",
			MusicalKeySignature: "sharp",
			MajorMinor:          "Major",
		},
	}

	fileTags := make(map[string][]string)
	fileTags[testAudioFiles[0].Title] = []string{"classical", "light", "instrumental"}
	fileTags[testAudioFiles[1].Title] = []string{"guitar", "grunge", "dark"}
	fileTags[testAudioFiles[2].Title] = []string{"guitar", "light", "slow"}
	fileTags[testAudioFiles[3].Title] = []string{"intense", "supertrap", "Gunna"}
	fileTags[testAudioFiles[4].Title] = []string{"sza", "vocals", "soul"}
	fileTags[testAudioFiles[5].Title] = []string{"sza", "funk"}
	fileTags[testAudioFiles[6].Title] = []string{"Hip Hop", "Strings", "Guitar"}
	fileTags[testAudioFiles[7].Title] = []string{"guitar", "vocals", "slow"}
	fileTags[testAudioFiles[6].Title] = []string{"Hip Hop", "Drake", "trap"}
	fileTags[testAudioFiles[6].Title] = []string{"instrumental", "kaytranada", "simple"}
	fileTags[testAudioFiles[2].Title] = []string{"grunge", "heavy", "electric guitar"}

	for _, file := range testAudioFiles {
		id, err := Query.AddTestAudioFile(ctx, file)
		util.EMsg(err, "adding test audio file")
		fmt.Printf("id:  %s\n", id)

		if file.Title == "Swifty" {
			err := Query.UpdateFeaturedTrack(ctx, id)
			if err != nil {
				return fmt.Errorf("error setting featured track: %s\n", err)
			}
		}

		if err != nil {
			return fmt.Errorf("error adding file %s: %s", file.Title, err)
		}
		for _, tag := range fileTags[file.Title] {
			tag_id, err := Query.GetTagId(ctx, tag)
			if err != nil {
				return fmt.Errorf("error getting tag Id: %s\n", err)
			}
			if tag_id == 0 {
				tag_id, err = Query.InsertTag(ctx, tag)
				if err != nil {
					return fmt.Errorf("error inserting tag Id: %s\n", err)
				}
			}

			err = Query.TestAddAudioFileTag(ctx, sqlc.TestAddAudioFileTagParams{
				TagID:       tag_id,
				AudioFileID: id,
			})

			if err != nil {
				return fmt.Errorf("error adding audio file tag: %s\n", err)
			}
		}
	}

	err = initPins(ctx, kevin.ID)
	if err != nil {
		return err
	}
	err = initPins(ctx, marty.ID)
	if err != nil {
		return err
	}

	followingMap := make(map[pgtype.UUID][]pgtype.UUID)
	followingMap[kevin.ID] = []pgtype.UUID{blake.ID, marty.ID}
	followingMap[marty.ID] = []pgtype.UUID{kevin.ID, john.ID}

	for follower, followingArr := range followingMap {
		for _, following := range followingArr {
			err = Query.TestFollowUser(ctx, sqlc.TestFollowUserParams{
				FollowerID:  follower,
				FollowingID: following,
			})
			if err != nil {
				return fmt.Errorf("error TestFollowUser: %s\n", err)
			}
		}
	}

	return nil
}

func initPins(ctx context.Context, userId pgtype.UUID) error {
	files, err := Query.GetFourUserAudioFiles(ctx, userId)
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		pin := sqlc.InsertTestPinParams{
			UserID: userId,
			FileID: files[i].ID,
		}
		err = Query.InsertTestPin(ctx, pin)
		if err != nil {
			return err
		}
	}
	return nil
}
