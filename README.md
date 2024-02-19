# Orhun: Türkçe Genel Kullanımlı Programlama Dili

Şu an için, sadece standard girdiden gelen programlar ile test edebilirsiniz.

```go build  ./...```

```cat deneme.orhun | ./orhun```

veya:

```cat deneme.orhun | go run main.go scan.go parse.go walk.go```

örnek bir orhun programı:
```
giriş:

	yeni tamsayı no = 2
	yeni tamsayı no2 = 3

	no = 14
	no'yu yazdır

	no'yu ve no2'yi yazdır

	yeni önerme test = yanlış

	test = (no <= 2)


	eğer test doğruysa:
		no = no + 3
	.
	değilse:
		no = no - 10
	.

	eğer no < 20 doğruysa yinele:
		no = no + 1
	.

	test = test veya doğru

	no'yu yazdır
.
```

## Yardım istenen birtakım mevzular

Test Programlarına ihtiyacımız var.

Genel anlamda, her türlü yardım hoşgörülür.

## İletişim

harmankaya@mshyazilim.com 

## Lisans

MIT Lisansı

## Hedefler:

#### v0.1:
globalde ve yerelde fonksiyon ve degisken tanimlari +

struct/ozne/yapi'ların eklenmesi. +

#### v0.2:

Go tarzı ":=" tanımları ile sözdiziminde basitleşmeye gitmek istiyorum.

karakter/byte/uint8, float ve dizi veri tiplerinin desteklenmesi.

#### v0.3:

UTF-8 stringleri

Modul desteği

#### Ötesi:
Sözdizimi ile ilgili fikirler oturduktan sonra, Orhunu derlenen bir dile dönüştürmek 
