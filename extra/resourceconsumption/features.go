package resourceconsumption

import (
	"errors"
	"github.com/cloudimpl/next-coder-sdk/polycode"
)

func CreateFeature(ctx polycode.ServiceContext, feature Feature) error {
	featuresCol := ctx.Db().Collection("polycode_Features")
	return featuresCol.InsertOne(feature)
}

func ConsumeFeature(ctx polycode.ServiceContext, id string, count int) (Feature, error) {
	featuresCol := ctx.Db().Collection("polycode_Features")

	var feature Feature
	exist, err := featuresCol.GetOne(id, &feature)
	if err != nil {
		return Feature{}, err
	}

	if !exist {
		return Feature{}, errors.New("feature not available")
	}

	consumed := float64(count) * feature.UnitCost
	if consumed > feature.Remaining {
		return Feature{}, errors.New("not enough credit to use feature")
	}

	feature.Used += consumed
	feature.Remaining -= consumed

	err = featuresCol.UpdateOne(feature)
	if err != nil {
		return Feature{}, err
	}

	return feature, nil
}
