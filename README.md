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

Fonksiyonların syntax i nasıl olmalı biraz düşünmem lazım.
Fikriniz varsa lütfen issue kısmında belirtmenizi rica ederim.

Test Programlarına ihtiyacımız var.

Genel anlamda, her türlü yardım hoşgörülür.

## İletişim

harmankaya@mshyazilim.com 

## Lisans

MIT Lisansı

## v0.1 Hedefleri:
globalde ve yerelde fonksiyon ve degisken tanimlari
struct/ozne/yapi tanimlamalari

sonraki hedefler?
Modul destegi ve syntaxi

LLVM 
