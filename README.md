
## Usage:

#### Insertion

Ici, nous avons une insertion simple dans une image. Le texte à inséré se trouve dans le fichier text.txt.

```
go-lsb -insert image_src.bmp image_dst.bmp text.txt
```

Si maintenant nous voulons chiffrer les données insérées il suffit de rajouter l'option key, comme ci-dessous.
```
go-lsb -key MYKEY -insert image_src.bmp image_dst.bmp text.txt
```

Enfin si nous voulons insérer les données de manières "aléatoire"

```
go-lsb -key MYKEY -seed THISISMYSEED -insert image_src.bmp image_dst.bmp text.txt
```

#### Récupération
Après avoir vu comment insérer les données nous allons voir comment les récupérer. Dans un premier temps, nous allons efféctuer une récupération simple.
```
go-lsb -retrive image_dst.bmp
```

Le contenu du fichier text.txt sera alors affiché sur la sortie standard.


Ensuite comme pour l'insertion, si vous voulez dechiffrer le text, il suffit de rajouter l'option key.
```
go-lsb -key MYKEY -retrive image_dst.bmp
```

Et enfin pour finir si vous voulez récupérer le contenu qui a été inséré aléatoirement, il suffit de rajouter l'option seed.
```
go-lsb -key MYKEY -seed THISISMYSEED -retrive image_dst.bmp
```

#### Détéction
Pour lancer la détéction sur une image, il suffit de taper la commande suivante
```
go-lsb -detect image_dst.bmp
```