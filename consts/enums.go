package consts

type PageEnum int

const (
	ExplorePage PageEnum = iota
	SearchPage
	MyDownloadsPage
	MyUploadsPage
)

type ComponentEnum int

const (
	TagBarComp ComponentEnum = iota
)

type Modal int

const (
	UploadModal Modal = iota
	CreateAccountModal
	LoginModal
	EditProfileModal
)

type SectionType int

const (
	CardCarouselSection SectionType = iota
	OtherUserPinsSection
	MyProfilePinsSection
)

type MusicalKey string

const (
	A MusicalKey = "A"
	B MusicalKey = "B"
	C MusicalKey = "C"
	D MusicalKey = "D"
	E MusicalKey = "E"
	F MusicalKey = "F"
	G MusicalKey = "G"
)

type MusicalKeySignature string

const (
	Flat    MusicalKeySignature = "flat"
	Natural MusicalKeySignature = "natural"
	Sharp   MusicalKeySignature = "sharp"
)

type MajorMinor string

const (
	Major MajorMinor = "Major"
	Minor MajorMinor = "Minor"
)

// type Head
