/*
1) Barrer una lista de más de 150 ítems ids en el servicio público:

https://api.mercadolibre.com/sites/MLA/search?q=chromecast&limit=50#json

2) Por cada resultado, realizar el correspondiente GET por Item_Id al recurso público:

https://api.mercadolibre.com/items/{Item_Id}
*/

package meli
