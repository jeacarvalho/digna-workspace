package exchange

import (
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type OfferStatus string

const (
	OfferActive    OfferStatus = "ACTIVE"
	OfferInactive  OfferStatus = "INACTIVE"
	OfferCompleted OfferStatus = "COMPLETED"
)

type Offer struct {
	ID          string
	EntityID    string
	EntityName  string
	Product     string
	Quantity    int64
	Unit        string
	Price       int64
	Currency    string
	Description string
	Status      OfferStatus
	CreatedAt   int64
	ExpiresAt   int64
}

type OfferRegistry struct {
	offers map[string][]Offer
}

func NewOfferRegistry() *OfferRegistry {
	return &OfferRegistry{
		offers: make(map[string][]Offer),
	}
}

func (or *OfferRegistry) Publish(offer Offer) error {
	if offer.EntityID == "" {
		return fmt.Errorf("entity_id is required")
	}
	if offer.Product == "" {
		return fmt.Errorf("product is required")
	}
	if offer.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	if offer.ID == "" {
		offer.ID = fmt.Sprintf("OFFER-%s-%d", offer.EntityID, time.Now().Unix())
	}

	if offer.CreatedAt == 0 {
		offer.CreatedAt = time.Now().Unix()
	}

	if offer.Status == "" {
		offer.Status = OfferActive
	}

	or.offers[offer.EntityID] = append(or.offers[offer.EntityID], offer)
	return nil
}

func (or *OfferRegistry) ListByEntity(entityID string) []Offer {
	if offers, ok := or.offers[entityID]; ok {
		return offers
	}
	return []Offer{}
}

func (or *OfferRegistry) ListAll() []Offer {
	var allOffers []Offer
	for _, offers := range or.offers {
		for _, offer := range offers {
			if offer.Status == OfferActive {
				allOffers = append(allOffers, offer)
			}
		}
	}
	return allOffers
}

func (or *OfferRegistry) ListByProduct(product string) []Offer {
	var matchingOffers []Offer
	for _, offers := range or.offers {
		for _, offer := range offers {
			if offer.Product == product && offer.Status == OfferActive {
				matchingOffers = append(matchingOffers, offer)
			}
		}
	}
	return matchingOffers
}

func (or *OfferRegistry) GetOffer(offerID string) (*Offer, error) {
	for _, offers := range or.offers {
		for i := range offers {
			if offers[i].ID == offerID {
				return &offers[i], nil
			}
		}
	}
	return nil, fmt.Errorf("offer not found: %s", offerID)
}

func (or *OfferRegistry) UpdateStatus(offerID string, status OfferStatus) error {
	for _, offers := range or.offers {
		for i := range offers {
			if offers[i].ID == offerID {
				offers[i].Status = status
				return nil
			}
		}
	}
	return fmt.Errorf("offer not found: %s", offerID)
}

func (or *OfferRegistry) CountActive() int {
	count := 0
	for _, offers := range or.offers {
		for _, offer := range offers {
			if offer.Status == OfferActive {
				count++
			}
		}
	}
	return count
}

type IntercoopService struct {
	registry         *OfferRegistry
	lifecycleManager lifecycle.LifecycleManager
}

func NewIntercoopService(lm lifecycle.LifecycleManager) *IntercoopService {
	return &IntercoopService{
		registry:         NewOfferRegistry(),
		lifecycleManager: lm,
	}
}

func (is *IntercoopService) PublishOffer(entityID string, product string, quantity int64, price int64, description string) (*Offer, error) {
	offer := Offer{
		EntityID:    entityID,
		Product:     product,
		Quantity:    quantity,
		Unit:        "unidade",
		Price:       price,
		Currency:    "BRL",
		Description: description,
		Status:      OfferActive,
	}

	if err := is.registry.Publish(offer); err != nil {
		return nil, err
	}

	published := is.registry.ListByEntity(entityID)
	if len(published) > 0 {
		return &published[len(published)-1], nil
	}

	return nil, fmt.Errorf("failed to publish offer")
}

func (is *IntercoopService) DiscoverOffers(productFilter string) []Offer {
	if productFilter != "" {
		return is.registry.ListByProduct(productFilter)
	}
	return is.registry.ListAll()
}

func (is *IntercoopService) GetEntityOffers(entityID string) []Offer {
	return is.registry.ListByEntity(entityID)
}
