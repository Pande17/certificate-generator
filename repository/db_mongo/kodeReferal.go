package dbmongo

type KodeReferal struct {
	ReferralID uint64 `json:"referral_id" bson:"referral_id"`
	Divisi     string `json:"divisi" bson:"divisi"`
	BulanRilis string `json:"bulan_rilis" bson:"bulan_rilis"`
	TahunRilis string `json:"tahun_rilis" bson:"tahun_rilis"`
}
