package sqlc

import (
	"fmt"
	"github.com/jmarren/deepfried/util"
	"net/url"
	"os"
)

func (p *Playable) AdjustSrcs() {
	if os.Getenv("ENV") == "prod" {
		audioEscaped := url.PathEscape(p.AudioSrc)
		fmt.Printf("\t\tAUDIO_SRC_ESCAPED: %s\n", audioEscaped)

		artworkSrc := fmt.Sprintf("https://%s/artwork/%s/%s", os.Getenv("STATIC_DOMAIN"), util.UuidString(p.UserID), p.ArtworkSrc)
		audioSrc := fmt.Sprintf("https://%s/audio/%s/%s", os.Getenv("STATIC_DOMAIN"), util.UuidString(p.UserID), audioEscaped)

		artworkSrcUrl, err := url.Parse(artworkSrc)
		util.EMsg(err, "parsing artworkSrc url")

		audioSrcUrl, err := url.Parse(audioSrc)
		util.EMsg(err, "parsing audioSrc url")

		fmt.Printf("\t\taudioSrcUrl.String(): %s \n", audioSrcUrl.String())
		// fmt.Printf("artworkSrcUrl.String(): %s ", audioSrcUrl.String())

		p.ArtworkSrc = artworkSrcUrl.String()
		p.AudioSrc = audioSrcUrl.String()

	} else {
		p.ArtworkSrc = fmt.Sprintf("%sartwork/%s/%s", os.Getenv("STATIC_DOMAIN"), util.UuidString(p.UserID), p.ArtworkSrc)
		p.AudioSrc = fmt.Sprintf("%saudio/%s", os.Getenv("STATIC_DOMAIN"), p.AudioSrc)
	}
}

func (u *User) GetProfilePhotoSrc() string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/profile_photos/%s", os.Getenv("STATIC_DOMAIN"), util.UuidString(u.ID))
	}
	return fmt.Sprintf("%sprofile_photos/%s", os.Getenv("STATIC_DOMAIN"), util.UuidString(u.ID))
}

/*
func (u *User) AmIFollowing(ctx context.Context, myId pgtype.UUID) bool {

// type IsUserFollowingUserParams struct {
// 	FollowerID  pgtype.UUID
// 	FollowingID pgtype.UUID
// }

	yesIam, err  := isUserFollowingUser(ctx, IsUserFollowingUserParams{
		FollowerID:
	})
}
*/
