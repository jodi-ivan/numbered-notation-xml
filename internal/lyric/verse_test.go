package lyric

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_lyricInteractor_RenderVerse(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	verse1a := `[ [{"word":"Suci,","breakdown":[{"text":"Su","type":"begin"},{"text":"ci,","type":"end"}]},{"word":"suci,","breakdown":[{"text":"su","type":"begin"},{"text":"ci,","type":"end"}]},{"word":"suci!","breakdown":[{"text":"su","type":"begin"},{"text":"ci!","type":"end"}]},{"word":"Kaum","breakdown":[{"text":"Kaum","type":"single","combine":true,"breakdown":[{"text":"K"},{"text":"au","underline":true},{"text":"m"}]}]},{"word":"kudus","breakdown":[{"text":"ku","type":"begin"},{"text":"dus","type":"end"}]},{"word":"tersungkur","breakdown":[{"text":"ter","type":"begin"},{"text":"sung","type":"middle"},{"text":"kur","type":"end"}]}], [{"word":"di","breakdown":[{"text":"di","type":"single"}]},{"word":"depan","breakdown":[{"text":"de","type":"begin"},{"text":"pan","type":"end"}]},{"word":"takhtaMu","breakdown":[{"text":"takh","type":"begin"},{"text":"ta","type":"middle"},{"text":"Mu","type":"end"}]},{"word":"memb'ri","breakdown":[{"text":"mem","type":"begin"},{"text":"b'ri","type":"end"}]},{"word":"mahkotanya.","breakdown":[{"text":"mah","type":"begin"},{"text":"ko","type":"middle"},{"text":"ta","type":"middle"},{"text":"nya.","type":"end"}]}], [{"word":"Segenap","breakdown":[{"text":"Se","type":"begin"},{"text":"ge","type":"middle"},{"text":"nap","type":"end"}]},{"word":"malaikai","breakdown":[{"text":"ma","type":"begin"},{"text":"lai","type":"middle","combine":true,"breakdown":[{"text":"l"},{"text":"ai","underline":true}]},{"text":"kat","type":"end"}]},{"word":"sujud","breakdown":[{"text":"su","type":"begin"},{"text":"jud","type":"end"}]},{"word":"menyembahMu,","breakdown":[{"text":"me","type":"begin"},{"text":"nyem","type":"middle"},{"text":"bah","type":"middle"},{"text":"Mu,","type":"end"}]}], [{"word":"Tuhan,","breakdown":[{"text":"Tu","type":"begin"},{"text":"han","type":"end"}]},{"word":"Yang","breakdown":[{"text":"Yang","type":"single"}]},{"word":"Ada","breakdown":[{"text":"A","type":"begin"},{"text":"da","type":"end"}]},{"word":"s'lama-lamanya","breakdown":[{"text":"s'la","type":"begin"},{"text":"ma","type":"middle"},{"text":"la","type":"middle"},{"text":"ma","type":"middle"},{"text":"nya","type":"end"}]}]]`

	_ = verse1a

	verse2a := `[[{"breakdown":[{"text":"Ke","type":"begin"},{"text":"ru","type":"middle"},{"text":"bim","type":"end"}],"word":"Kerubim"},{"breakdown":[{"text":"dan","type":"single"}],"word":"dan"},{"breakdown":[{"text":"se","type":"begin"},{"text":"ra","type":"middle"},{"text":"fim","type":"end"}],"word":"serafim"}],[{"breakdown":[{"text":"me","type":"begin"},{"text":"mu","type":"middle"},{"text":"lia","combine":true,"type":"middle","breakdown":[{"text":"l"},{"text":"ia","underline":true}]},{"text":"kan","type":"end"}],"word":"memuliakan"},{"breakdown":[{"text":"Yang","type":"single"}],"word":"Yang"},{"breakdown":[{"text":"Tri","type":"begin"},{"text":"su","type":"middle"},{"text":"ci;","type":"end"}],"word":"Trisuci;"}],[{"breakdown":[{"text":"pa","type":"begin"},{"text":"ra","type":"end"}],"word":"para"},{"breakdown":[{"text":"ra","type":"begin"},{"text":"sul","type":"end"}],"word":"rasul"},{"breakdown":[{"text":"dan","type":"single"}],"word":"dan"},{"breakdown":[{"text":"na","type":"begin"},{"text":"bi,","type":"end"}],"word":"nabi,"}],[{"breakdown":[{"text":"mar","type":"begin"},{"text":"tir","type":"end"}],"word":"martir"},{"breakdown":[{"text":"yang","type":"single"}],"word":"yang"},{"breakdown":[{"text":"ber","type":"begin"},{"text":"ju","type":"middle"},{"text":"bah","type":"end"}],"word":"berjubah"},{"breakdown":[{"text":"pu","type":"begin"},{"text":"tih","type":"end"}],"word":"putih"}],[{"breakdown":[{"text":"G're","type":"begin"},{"text":"ja","type":"end"}],"word":"G'reja"},{"breakdown":[{"text":"yang","type":"single"}],"word":"yang"},{"breakdown":[{"text":"ku","type":"begin"},{"text":"dus,","type":"end"}],"word":"kudus,"},{"breakdown":[{"text":"e","type":"begin"},{"text":"sa,","type":"end"}],"word":"esa,"}],[{"breakdown":[{"text":"ke","type":"begin"},{"text":"pa","type":"middle"},{"text":"da","type":"middle"},{"text":"Mu","type":"end"}],"word":"kepadaMu"},{"breakdown":[{"text":"me","type":"begin"},{"text":"nyem","type":"middle"},{"text":"bah.","type":"end"}],"word":"menyembah."}]]`
	verse2b := `[[{"breakdown":[{"text":"Takh","type":"begin"},{"text":"ta","type":"middle"},{"text":"Mu","type":"end"}],"word":"TakhtaMu"},{"breakdown":[{"text":"ke","type":"begin"},{"text":"kal","type":"end"}],"word":"kekal"},{"breakdown":[{"text":"te","type":"begin"},{"text":"guh","type":"end"}],"word":"teguh"}],[{"breakdown":[{"text":"pa","type":"begin"},{"text":"da","type":"end"}],"word":"pada"},{"breakdown":[{"text":"Si","type":"begin"},{"text":"si","type":"end"}],"word":"Sisi"},{"breakdown":[{"text":"ka","type":"begin"},{"text":"nan","type":"end"}],"word":"kanan"},{"breakdown":[{"text":"Ba","type":"begin"},{"text":"pa;","type":"end"}],"word":"Bapa;"}],[{"breakdown":[{"text":"da","type":"begin"},{"text":"lam","type":"end"}],"word":"dalam"},{"breakdown":[{"text":"peng","type":"begin"},{"text":"ha","type":"middle"},{"text":"ki","type":"middle"},{"text":"man","type":"middle"},{"text":"Mu,","type":"end"}],"word":"penghakimanMu,"}],[{"breakdown":[{"text":"to","type":"begin"},{"text":"long","type":"end"}],"word":"tolong"},{"breakdown":[{"text":"u","type":"begin"},{"text":"mat","type":"middle"},{"text":"Mu","type":"end"}],"word":"umatMu"},{"breakdown":[{"text":"yang","type":"single"}],"word":"yang"},{"breakdown":[{"text":"pa","type":"begin"},{"text":"pa:","type":"end"}],"word":"papa:"}],[{"breakdown":[{"text":"di","type":"begin"},{"text":"ri","type":"end"}],"word":"diri"},{"breakdown":[{"text":"ka","type":"begin"},{"text":"mi","type":"end"}],"word":"kami"},{"breakdown":[{"text":"yang","type":"single"}],"word":"yang"},{"breakdown":[{"text":"le","type":"begin"},{"text":"mah","type":"end"}],"word":"lemah"}],[{"breakdown":[{"text":"da","type":"begin"},{"text":"lam","type":"end"}],"word":"dalam"},{"breakdown":[{"text":"Di","type":"begin"},{"text":"kau","type":"end"}],"word":"Dikau"},{"breakdown":[{"text":"s'la","type":"begin"},{"text":"mat","type":"middle"},{"text":"lah!","type":"end"}],"word":"s'lamatlah!"}]]`

	_ = verse2a
	_ = verse2b

	verse3a := `[[{"breakdown":[{"text":"Di","type":"single"}],"word":"Di"},{"breakdown":[{"text":"ha","type":"begin"},{"text":"da","type":"middle"},{"text":"pan","type":"end"}],"word":"hadapan"},{"breakdown":[{"text":"tak","type":"begin"},{"text":"ta","type":"end"}],"word":"takta"},{"breakdown":[{"text":"rah","type":"begin"},{"text":"mat","type":"end"}],"word":"rahmat"}],[{"breakdown":[{"text":"a","type":"begin"},{"text":"ku","type":"end"}],"word":"aku"},{"breakdown":[{"text":"me","type":"begin"},{"text":"nyem","type":"middle"},{"text":"bah,","type":"end"}],"word":"menyembah,"}],[{"breakdown":[{"text":"tun","type":"begin"},{"text":"duk","type":"end"}],"word":"tunduk"},{"breakdown":[{"text":"da","type":"begin"},{"text":"lam","type":"end"}],"word":"dalam"},{"breakdown":[{"text":"pe","type":"begin"},{"text":"nye","type":"middle"},{"text":"sa","type":"middle"},{"text":"lan","type":"end"}],"word":"penyesalan"}],[{"breakdown":[{"text":"Tu","type":"begin"},{"text":"han","type":"end"}],"word":"Tuhan"},{"breakdown":[{"text":"to","type":"begin"},{"text":"long","type":"middle"},{"text":"lah!","type":"end"}],"word":"tolonglah!"}]]`
	verse3b := `[[{"breakdown":[{"text":"\"I","type":"begin"},{"text":"ni","type":"end"}],"word":"\"Ini"},{"breakdown":[{"text":"sa","type":"begin"},{"text":"ja","type":"end"}],"word":"saja"},{"breakdown":[{"text":"an","type":"begin"},{"text":"da","type":"middle"},{"text":"lan","type":"middle"},{"text":"ku:","type":"end"}],"word":"andalanku:"}],[{"breakdown":[{"text":"ja","type":"begin"},{"text":"sa","type":"end"}],"word":"jasa"},{"breakdown":[{"text":"kur","type":"begin"},{"text":"ban","type":"middle"},{"text":"Mu.","type":"end"}],"word":"kurbanMu."}],[{"breakdown":[{"text":"Ha","type":"begin"},{"text":"ti","type":"middle"},{"text":"ku","type":"end"}],"word":"Hatiku"},{"breakdown":[{"text":"yang","type":"single"}],"word":"yang"},{"breakdown":[{"text":"han","type":"begin"},{"text":"cur","type":"end"}],"word":"hancur"},{"breakdown":[{"text":"lu","type":"begin"},{"text":"luh","type":"end"}],"word":"luluh"}],[{"breakdown":[{"text":"bu","type":"begin"},{"text":"at","type":"middle"},{"text":"lah","type":"end"}],"word":"buatlah"},{"breakdown":[{"text":"sem","type":"begin"},{"text":"buh.","type":"end"}],"word":"sembuh."}]]`
	verse3c := `[[{"breakdown":[{"text":"Ka","type":"begin"},{"text":"u","type":"middle"},{"text":"lah","type":"end"}],"word":"Kaulah"},{"breakdown":[{"text":"Sum","type":"begin"},{"text":"ber","type":"end"}],"word":"Sumber"},{"breakdown":[{"text":"peng","type":"begin"},{"text":"hi","type":"middle"},{"text":"bu","type":"middle"},{"text":"ran","type":"end"}],"word":"penghiburan"}],[{"breakdown":[{"text":"Ra","type":"begin"},{"text":"ja","type":"end"}],"word":"Raja"},{"breakdown":[{"text":"hi","type":"begin"},{"text":"dup","type":"middle"},{"text":"ku.","type":"end"}],"word":"hidupku."}],[{"breakdown":[{"text":"Baik","combine":true,"type":"single","breakdown":[{"text":"B"},{"text":"ai","underline":true},{"text":"k"}]}],"word":"Baik"},{"breakdown":[{"text":"di","type":"single"}],"word":"di"},{"breakdown":[{"text":"bu","type":"begin"},{"text":"mi","type":"end"}],"word":"bumi"},{"breakdown":[{"text":"baik","combine":true,"type":"single","breakdown":[{"text":"b"},{"text":"ai","underline":true},{"text":"k"}]}],"word":"baik"},{"breakdown":[{"text":"di","type":"single"}],"word":"di"},{"breakdown":[{"text":"sor","type":"begin"},{"text":"ga,","type":"end"}],"word":"sorga,"}],[{"breakdown":[{"text":"sia","combine":true,"type":"begin","breakdown":[{"text":"sia","underline":true}]},{"text":"pa","type":"end"}],"word":"siapa"},{"breakdown":[{"text":"ban","type":"begin"},{"text":"ding","type":"middle"},{"text":"Mu?","type":"end"}],"word":"bandingMu?"}]]`

	type args struct {
		y      int
		verses []repository.HymnVerse
	}
	tests := []struct {
		name       string
		args       args
		want       VerseInfo
		initCanvas func(*gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "KJ-002. 1column",
			args: args{
				y: 100,
				verses: []repository.HymnVerse{
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 2, Valid: true},
						VerseNum: sql.NullInt32{Int32: 2, Valid: true},
						Row:      sql.NullInt16{Int16: 1, Valid: true},
						StyleRow: sql.NullInt32{Int32: 12, Valid: true},
						Col:      sql.NullInt16{Int16: 1, Valid: true},
						Content:  sql.NullString{String: verse1a, Valid: true},
					},
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='verses'", "style='font-family:Caladea'")
				canv.EXPECT().Group("class='verse'")
				canv.EXPECT().Group()

				canv.EXPECT().Text(341, 100, "2. ")
				canv.EXPECT().Text(360, 100, " Suci, suci, suci! Kaum kudus tersungkur")
				canv.EXPECT().Text(360, 125, " di depan takhtaMu memb'ri mahkotanya.")
				canv.EXPECT().Text(360, 150, " Segenap malaikai sujud menyembahMu,")
				canv.EXPECT().Text(360, 175, " Tuhan, Yang Ada s'lama-lamanya")
				canv.EXPECT().Gend().Times(3)

				canv.EXPECT().Qbez(472, 102, 480, 107, 489, 102, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Qbez(442, 152, 447, 157, 454, 152, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				return canv
			},
			want: VerseInfo{
				MarginBottom: 235,
			},
		},
		{
			name: "KJ-005. 2column",
			args: args{
				y: 100,
				verses: []repository.HymnVerse{
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 5, Valid: true},
						VerseNum: sql.NullInt32{Int32: 2, Valid: true},
						Row:      sql.NullInt16{Int16: 1, Valid: true},
						StyleRow: sql.NullInt32{Int32: 6, Valid: true},
						Col:      sql.NullInt16{Int16: 1, Valid: true},
						Content:  sql.NullString{String: verse2a, Valid: true},
					},
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 5, Valid: true},
						VerseNum: sql.NullInt32{Int32: 3, Valid: true},
						Row:      sql.NullInt16{Int16: 1, Valid: true},
						StyleRow: sql.NullInt32{Int32: 6, Valid: true},
						Col:      sql.NullInt16{Int16: 2, Valid: true},
						Content:  sql.NullString{String: verse2b, Valid: true},
					},
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='verses'", "style='font-family:Caladea'")
				canv.EXPECT().Group("class='verse'").Times(2)
				canv.EXPECT().Group().Times(2)

				canv.EXPECT().Text(155, 100, "2. ")
				canv.EXPECT().Text(174, 100, " Kerubim dan serafim")
				canv.EXPECT().Text(174, 125, " memuliakan Yang Trisuci;")
				canv.EXPECT().Text(174, 150, " para rasul dan nabi,")
				canv.EXPECT().Text(174, 175, " martir yang berjubah putih")
				canv.EXPECT().Text(174, 200, " G'reja yang kudus, esa,")
				canv.EXPECT().Text(174, 225, " kepadaMu menyembah.")

				canv.EXPECT().Gend().Times(5)

				canv.EXPECT().Qbez(220, 127, 225, 132, 232, 127, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")

				canv.EXPECT().Text(391, 100, "3. ")
				canv.EXPECT().Text(410, 100, " TakhtaMu kekal teguh")
				canv.EXPECT().Text(410, 125, " pada Sisi kanan Bapa;")
				canv.EXPECT().Text(410, 150, " dalam penghakimanMu,")
				canv.EXPECT().Text(410, 175, " tolong umatMu yang papa:")
				canv.EXPECT().Text(410, 200, " diri kami yang lemah")
				canv.EXPECT().Text(410, 225, " dalam Dikau s'lamatlah!")

				return canv
			},
			want: VerseInfo{
				MarginBottom: 285,
			},
		},
		{
			name: "KJ-026. mixed 2 and 1 column",
			args: args{
				y: 100,
				verses: []repository.HymnVerse{
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 5, Valid: true},
						VerseNum: sql.NullInt32{Int32: 2, Valid: true},
						Row:      sql.NullInt16{Int16: 1, Valid: true},
						StyleRow: sql.NullInt32{Int32: 6, Valid: true},
						Col:      sql.NullInt16{Int16: 1, Valid: true},
						Content:  sql.NullString{String: verse3a, Valid: true},
					},
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 5, Valid: true},
						VerseNum: sql.NullInt32{Int32: 3, Valid: true},
						Row:      sql.NullInt16{Int16: 1, Valid: true},
						StyleRow: sql.NullInt32{Int32: 6, Valid: true},
						Col:      sql.NullInt16{Int16: 2, Valid: true},
						Content:  sql.NullString{String: verse3b, Valid: true},
					},
					repository.HymnVerse{
						Number:   sql.NullInt32{Int32: 5, Valid: true},
						VerseNum: sql.NullInt32{Int32: 4, Valid: true},
						Row:      sql.NullInt16{Int16: 2, Valid: true},
						Col:      sql.NullInt16{Int16: 1, Valid: true},
						Content:  sql.NullString{String: verse3c, Valid: true},
					},
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='verses'", "style='font-family:Caladea'")
				canv.EXPECT().Group("class='verse'").Times(3)
				canv.EXPECT().Group().Times(3)

				canv.EXPECT().Text(157, 100, "2. ")
				canv.EXPECT().Text(176, 100, " Di hadapan takta rahmat")
				canv.EXPECT().Text(176, 125, " aku menyembah,")
				canv.EXPECT().Text(176, 150, " tunduk dalam penyesalan")
				canv.EXPECT().Text(176, 175, " Tuhan tolonglah!")

				canv.EXPECT().Gend().Times(7)

				canv.EXPECT().Text(401, 100, "3. ")
				canv.EXPECT().Text(420, 100, " \"Ini saja andalanku:")
				canv.EXPECT().Text(420, 125, " jasa kurbanMu.")
				canv.EXPECT().Text(420, 150, " Hatiku yang hancur luluh")
				canv.EXPECT().Text(420, 175, " buatlah sembuh.")

				canv.EXPECT().Text(340, 235, "4. ")
				canv.EXPECT().Text(360, 235, " Kaulah Sumber penghiburan")
				canv.EXPECT().Text(360, 260, " Raja hidupku.")
				canv.EXPECT().Text(360, 285, " Baik di bumi baik di sorga,")
				canv.EXPECT().Text(360, 310, " siapa bandingMu?")

				canv.EXPECT().Qbez(369, 287, 374, 292, 381, 287, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Qbez(457, 287, 462, 292, 469, 287, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Qbez(360, 312, 369, 317, 378, 312, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")

				return canv
			},
			want: VerseInfo{
				MarginBottom: 370,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := lyricInteractor{}
			if got := li.RenderVerse(context.Background(), tt.initCanvas(ctrl), tt.args.y, tt.args.verses); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lyricInteractor.RenderVerse() = %v, want %v", got, tt.want)
			}
		})
	}
}
