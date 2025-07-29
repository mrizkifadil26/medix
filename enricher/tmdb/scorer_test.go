package tmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/mrizkifadil26/medix/enricher/tmdb"
)

const sampleFlowJSON = `{
  "page": 1,
  "results": [
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        16
      ],
      "id": 1281775,
      "original_language": "en",
      "original_title": "FLOW",
      "overview": "As the ‘flow’ progresses, the world where our ‘giant’ once stood gradually fades away, replaced by their inner worlds. Our ‘giant’ continues to advance endlessly, questioning the relationship between the world and oneself.",
      "popularity": 0.4784,
      "poster_path": null,
      "release_date": "2024-02-01",
      "title": "FLOW",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [],
      "id": 1283979,
      "original_language": "no",
      "original_title": "Flow",
      "overview": "Flow by Anne Haugsgjerd masterfully depicts the ephemeral beauty of life through a filmmaker at a crossroads, reflecting on aging, family and the flow of life. Haugsgjerd's blend of humor and euphoria, paired with her sensitive and imaginative direction, inspires a profound reflection on existence and elegantly poses the question of choosing a quieter life amidst uncertainty.",
      "popularity": 0.2219,
      "poster_path": null,
      "release_date": "2024-05-03",
      "title": "Flow",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        16
      ],
      "id": 1390255,
      "original_language": "en",
      "original_title": "Flow",
      "overview": "Lens-based video and photography combine with digital geometric shapes to imitate the invisible flow in nature.",
      "popularity": 0.1642,
      "poster_path": null,
      "release_date": "2024-11-09",
      "title": "Flow",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/pZfpUMEJzw9BZHtXuZzzOrRxtMa.jpg",
      "genre_ids": [
        16,
        12,
        14,
        10751
      ],
      "id": 823219,
      "original_language": "lv",
      "original_title": "Straume",
      "overview": "A solitary cat, displaced by a great flood, finds refuge on a boat with various species and must navigate the challenges of adapting to a transformed world together.",
      "popularity": 25.023,
      "poster_path": "/imKSymKBK7o73sajciEmndJoVkR.jpg",
      "release_date": "2024-08-29",
      "title": "Flow",
      "video": false,
      "vote_average": 8.2,
      "vote_count": 2151
    },
    {
      "adult": false,
      "backdrop_path": "/kO7ItfYphDZxfjjOxyWX0VUpGO4.jpg",
      "genre_ids": [
        28,
        10749,
        14
      ],
      "id": 1117098,
      "original_language": "zh",
      "original_title": "花千骨",
      "overview": "In order to avoid the resurrection of the demon, the immortal Bai Zihua led the immortal world to try to seal the demon again. He and his disciple, the orphan Hua Qiangu, develop a relationship, but face an even greater world crisis.",
      "popularity": 0.3537,
      "poster_path": "/dX3j2TY55mmUNxP2aFpucq03Ts2.jpg",
      "release_date": "2024-01-20",
      "title": "The Journey of Flower",
      "video": false,
      "vote_average": 4.5,
      "vote_count": 2
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [],
      "id": 1193987,
      "original_language": "ja",
      "original_title": "FUKUYAMA MASAHARU WE'RE BROS. TOUR 2024",
      "overview": "",
      "popularity": 0.0143,
      "poster_path": "/kKJeYMjDalLg9rSdS6pkLxHJCmX.jpg",
      "release_date": "2024-04-27",
      "title": "WE’RE BROS. TOUR 2024 Flowers and Bees, Tears and Music.",
      "video": false,
      "vote_average": 7.5,
      "vote_count": 1
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        18
      ],
      "id": 1306456,
      "original_language": "ja",
      "original_title": "雨花蓮歌",
      "overview": "Harumi is an ordinary university student, preoccupied both with relationships within her family and off-hand remarks by her friends. Her older sister Reiko wants to get married, but faces opposition from her mother and those around her. Gradually the two find themselves in constant conflict—the kind that develops between people who care about one another.",
      "popularity": 0.0071,
      "poster_path": "/x45nWhtuUafFgBoFkARHO7iqCDj.jpg",
      "release_date": "2024-07-14",
      "title": "Poems of Flower Rain",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/oMvSKafFfziHdw6qtChVVGDWMXS.jpg",
      "genre_ids": [
        99
      ],
      "id": 1274019,
      "original_language": "es",
      "original_title": "Ni siquiera las flores",
      "overview": "Mariana Viñoles has positioned her camera in the window. As lockdown puts everything on hold, little everyday dramas continue to play out in her street, as many secret tales waiting to be told. Off-screen, everyday dialogues outline the contours of domestic life. We are invited to prolong their existence, and marvel at the magical power of life.",
      "popularity": 0.0214,
      "poster_path": "/pqAcMoSfkINVISAvmn49KyQZDOQ.jpg",
      "release_date": "2024-04-16",
      "title": "Not Even the Flowers",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        18,
        36,
        10402
      ],
      "id": 1207442,
      "original_language": "pl",
      "original_title": "Idź pod prąd",
      "overview": "Second half of the 1970s. A few teenagers from the town of Ustrzyki Dolne, led by a charismatic and undisciplined student of the local school, Siczka, decided to become punks and get into punk rock and start a band called KSU.",
      "popularity": 0.2056,
      "poster_path": "/82NQ1ClmTRQ4BpCZcQWA3TLXFNH.jpg",
      "release_date": "2024-09-27",
      "title": "Go Against the Flow",
      "video": false,
      "vote_average": 5.667,
      "vote_count": 3
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        99
      ],
      "id": 1229789,
      "original_language": "en",
      "original_title": "Flowing Air",
      "overview": "Lane Lamoreaux was seriously injured in a paragliding accident and over a five year period he struggles to get his old life back again. Eventually Lane's forced to come face to face with his situation and the biggest decision of his life.",
      "popularity": 0,
      "poster_path": "/yqTNsM15HRGPQJmQon2KNFJT1NC.jpg",
      "release_date": "2024-01-11",
      "title": "Flowing Air",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/8hppDxJptuEx2MMlO18QZ87EDPy.jpg",
      "genre_ids": [
        99
      ],
      "id": 1262611,
      "original_language": "pt",
      "original_title": "Sem TÍtulo #9: Nem Todas as Flores da Falta",
      "overview": "Coleridge's flower in the garden of the paths that fork from Borges. Paper and film flowers; fossil flowers. \"Panorama of all the flowers of speech\" (Joyce), a \"flower full of the real, of the current\" (Wallace Stevens).",
      "popularity": 0.0566,
      "poster_path": "/oBO2lDyk4tVTKSmsZTH4YeDpImv.jpg",
      "release_date": "2024-04-05",
      "title": "Untitled #9: Nor All Flowers of Foul",
      "video": false,
      "vote_average": 5,
      "vote_count": 1
    },
    {
      "adult": false,
      "backdrop_path": "/5IZBx4UQpihV2Z3mfaOhyvAMcu8.jpg",
      "genre_ids": [
        16,
        10402
      ],
      "id": 1377493,
      "original_language": "ja",
      "original_title": "さくらみこ1st Live “flower fantasista!”",
      "overview": "Sakura Miko's 1st solo live concert, held at ARIAKE ARENA on October 26, 2024.",
      "popularity": 0.0214,
      "poster_path": "/sR0nXuHowmYWvowyezbyBE1wszB.jpg",
      "release_date": "2024-10-26",
      "title": "Sakura Miko 1st Live \"flower fantasista!\"",
      "video": true,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/8zFLbTHwgG3IXkvJnEZsinrMjXl.jpg",
      "genre_ids": [
        18
      ],
      "id": 1250893,
      "original_language": "id",
      "original_title": "Moving Flower",
      "overview": "Lily (8) decided to go to her father's house. Before she and her mom had to move out of town with her new family. Lily and Hendra (36) spent their time together that day.",
      "popularity": 0.0286,
      "poster_path": "/xhcWstFhQ7WXuN1CqmGTHOSbYAN.jpg",
      "release_date": "2024-02-27",
      "title": "Moving Flower",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        99
      ],
      "id": 1370628,
      "original_language": "en",
      "original_title": "la Flor del Camino",
      "overview": "In a few brief moments in the life of a young girl, we see her drawing on the window with a marker as she sings and chats, and walking along the main road on her way to meet her friend, meanwhile surrounded by cars and trucks racing dangerously close to her.",
      "popularity": 0.0071,
      "poster_path": "/xWxKCe6GeVZ4SrXV1E0qkU5fIlJ.jpg",
      "release_date": "2024-11-14",
      "title": "the Flower by the Road",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/8bMgA1vsfSXzVPxNeBUJxrti6Y2.jpg",
      "genre_ids": [
        18,
        9648
      ],
      "id": 1322087,
      "original_language": "ru",
      "original_title": "В потоке трех стихий",
      "overview": "Diving diver Vlad Fokin, who is going through a career crisis, receives an unexpected offer. He will have to ingratiate himself with a group of cliff divers suspected of a string of robberies. Perhaps this gang is connected with Vlad's former best friend, whose relationship ended six years ago. New experiences, thirst for risk and extreme heights consume the guy with his head.",
      "popularity": 0.106,
      "poster_path": "/h7QS39FzB0SgrMP9RmmDwLvMmsp.jpg",
      "release_date": "2024-07-25",
      "title": "In the Flow of the Three Elements",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [],
      "id": 1368674,
      "original_language": "en",
      "original_title": "Deep Dreaming Flowers to AltAir V.2",
      "overview": "Deep Dreaming Flowers is a recording of a live signal analog video/audio synthesizer performance with a voiced narration made in collaboration with an AI program. A speculative machine-guided psychedelic broadcast of an astral floral projection. A fictional telepathic transmission that saturates the boundaries of perception with interlacing signals of interconnected consciousness.",
      "popularity": 0,
      "poster_path": null,
      "release_date": "2024-10-19",
      "title": "Deep Dreaming Flowers to AltAir V.2",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": "/2A2Ckh2RqPhu4KSexCga0ARqJxO.jpg",
      "genre_ids": [
        18,
        10749,
        10402
      ],
      "id": 1252569,
      "original_language": "vi",
      "original_title": "Một Bông Hoa Mong manh",
      "overview": "Thach Thao is a beautiful girl with an angelic beautiful voice, and the first time she sang on stage at Itango Club, Thao immediately caught the eye of a music producer named Son. Then, she was invited to cooperate in filming a music video called \"A Fragile Flower\". After the music video was released, Thach Thao became a renowned singer. She is beautiful, talented, so many men love her. Then, which man will Thao's heart belong to? And what is the dangerous disease that she encountered?",
      "popularity": 0.1748,
      "poster_path": "/egR3nykprnHV4etxNFgHvJS3x9b.jpg",
      "release_date": "2024-03-29",
      "title": "A Fragile Flower",
      "video": false,
      "vote_average": 7,
      "vote_count": 1
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [],
      "id": 1368669,
      "original_language": "en",
      "original_title": "The Flower Cult of Amelia Earhart",
      "overview": "A synaptic celluloid requiem, propelling the High Priestess Aviator Earhart through far-sighted passages of flora, fauna, air, fire and water.",
      "popularity": 0.0453,
      "poster_path": "/wdOUX5feFrh83jUe0ykZhPOC92Q.jpg",
      "release_date": "2024-10-19",
      "title": "The Flower Cult of Amelia Earhart",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [],
      "id": 1368681,
      "original_language": "en",
      "original_title": "Deep Dreaming Flowers to AltAir V.2",
      "overview": "Deep Dreaming Flowers is a recording of a live signal analog video/audio synthesizer performance with a voiced narration made in collaboration with an AI program. A speculative machine-guided psychedelic broadcast of an astral floral projection. A fictional telepathic transmission that saturates the boundaries of perception with interlacing signals of interconnected consciousness.",
      "popularity": 0,
      "poster_path": null,
      "release_date": "2024-10-19",
      "title": "Deep Dreaming Flowers to AltAir V.2",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    },
    {
      "adult": false,
      "backdrop_path": null,
      "genre_ids": [
        16,
        14
      ],
      "id": 1307088,
      "original_language": "ja",
      "original_title": "ルカと太陽の花",
      "overview": "The story, set in a village covered in clouds, depicts the bond between a boy, Luka, and his childhood friend, Rosa. One day, Rosa, who has the power of the sun, hears from a mysterious duo that there is a way to clear the clouds. Luka, on the other hand, is attracted to a traveling band that visits the village and decides to leave...",
      "popularity": 0.0143,
      "poster_path": "/mjra06C8ivfdOoBhKLvWEhZ2jgO.jpg",
      "release_date": "2024-03-22",
      "title": "Luka and the Flower of the Sun",
      "video": false,
      "vote_average": 0,
      "vote_count": 0
    }
  ],
  "total_pages": 6,
  "total_results": 107
}`

func TestPickBestMovieMatch_Flow2024(t *testing.T) {
	var result tmdb.SearchResult
	if err := json.Unmarshal([]byte(sampleFlowJSON), &result); err != nil {
		t.Fatalf("failed to parse sample JSON: %v", err)
	}

	best := tmdb.PickBestMovieMatch(result.Results, "Flow", 2024)
	if best == nil {
		t.Fatal("no best match found")
	}

	if best.OriginalTitle != "Straume" {
		t.Errorf("expected 'Straume' as best match, got: %s (ID: %d)", best.OriginalTitle, best.ID)
	} else {
		t.Logf("✅ Best match: %s (%s), original: %s", best.Title, best.ReleaseDate, best.OriginalTitle)
	}
}
