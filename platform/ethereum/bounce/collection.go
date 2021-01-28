package bounce

import (
	"strconv"

	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/types"
)

var (
	nftVersion = "3.0" // opensea nft_version compatible
)

func (c *Client) GetCollections(owner string, coinIndex uint) (types.CollectionPage, error) {
	collections, err := c.getCollections(owner)
	if err != nil {
		return nil, err
	}
	return c.NormalizeCollections(collections, coinIndex, owner)

}

func (c *Client) GetCollectibles(owner, collectionID string, coinIndex uint) (types.CollectiblePage, error) {
	collectibles, err := c.getCollectibles(owner, collectionID)
	if err != nil {
		return nil, err
	}
	return c.NormalizeCollectibles(collectibles, coinIndex)
}

func (c *Client) NormalizeCollections(collections []Collection, coinIndex uint, owner string) (types.CollectionPage, error) {
	page := make(types.CollectionPage, 0)
	category := map[string]bool{}
	for _, cl := range collections {

		// existed category
		if _, ok := category[cl.ContractAddr]; ok {
			continue
		}

		total, err := strconv.Atoi(cl.Balance)
		if err != nil {
			continue
		}
		// skip empty info token
		if len(cl.TokenURI) == 0 {
			continue
		}

		info, err := fetchTokenURI(cl.TokenURI)
		if err != nil {
			return nil, err
		}

		// skip empty name/image
		if len(info.Name) == 0 || len(info.Image) == 0 {
			continue
		}

		page = append(page, types.Collection{
			Id:           cl.ContractAddr,
			Name:         info.Name,
			ImageUrl:     normalizeUrl(info.Image),
			Description:  info.Description,
			ExternalLink: normalizeUrl(cl.TokenURI),
			Total:        total,
			Address:      owner,
			Coin:         coinIndex,
			Type:         "ERC" + cl.TokenType,
		})
		category[cl.ContractAddr] = true
	}
	return page, nil
}

func (c *Client) NormalizeCollectibles(collectibles []Collectible, coinIndex uint) (types.CollectiblePage, error) {
	if len(collectibles) == 0 {
		return types.CollectiblePage{}, nil
	}
	page := make(types.CollectiblePage, 0)
	info, err := fetchTokenURI(collectibles[0].TokenURI)
	if err != nil {
		return nil, err
	}
	for _, c := range collectibles {
		category := c.ContractName
		if len(category) == 0 {
			category = info.Name
		}
		page = append(page, types.Collectible{
			ID:              blockatlas.GenCollectibleId(c.ContractAddr, c.TokenID),
			CollectionID:    c.ContractAddr,
			TokenID:         c.TokenID,
			ContractAddress: c.ContractAddr,
			Category:        category,
			ImageUrl:        normalizeUrl(info.Image),
			ExternalLink:    normalizeUrl(c.TokenURI),
			Type:            string(types.ERC721),
			Description:     info.Description,
			Coin:            coinIndex,
			Name:            info.Name,
			Version:         nftVersion,
		})
	}
	return page, nil
}
