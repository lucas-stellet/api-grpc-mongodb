package services

import (
	"context"
	"encoding/json"
	"fmt"
	"lucas-stellet/api-grpc-mongodb/pkg/db"
	"lucas-stellet/api-grpc-mongodb/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Applications struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	FolderName string             `json:"folderName" bson:"folderName"`
}

type UniversalApps struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	FolderName string             `json:"folderName" bson:"folderName"`
}

type Functions struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}
type Dependencies struct {
	Applications  []Applications  `json:"applications" bson:"applications"`
	Functions     []Functions     `json:"functions" bson:"functions"`
	UniversalApps []UniversalApps `json:"universalApps" bson:"universalApps"`
}

type ProjectMiddleware struct {
	Repository            string       `json:"repository" bson:"repository"`
	LogoURL               string       `json:"logoUrl" bson:"logoUrl"`
	Name                  string       `json:"name" bson:"name"`
	GatewayProjectVersion string       `json:"gatewayProjectVersion" bson:"gatewayProjectVersion"`
	IsActive              bool         `json:"isActive" bson:"isActive"`
	Description           string       `json:"description" bson:"description"`
	AutoTracking          []string     `json:"autoTracking" bson:"autoTracking"`
	SubscriberDatabase    string       `json:"subscriberDatabase" bson:"subscriberDatabase"`
	SubscriberRepository  string       `json:"subscriberRepository" bson:"subscriberRepository"`
	SubscriberName        string       `json:"subscriberName" bson:"subscriberName"`
	Dependencies          Dependencies `json:"dependencies" bson:"dependencies"`
}

func (p *ProjectMiddleware) GetData(projectID string) (string, error) {
	var projects []ProjectMiddleware

	_id, _ := primitive.ObjectIDFromHex(projectID)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":     _id,
				"deleted": false,
				"type":    "EXPOSE_API",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "subscribers",
				"localField":   "subscriber",
				"foreignField": "_id",
				"as":           "subscriber",
			},
		},
		{
			"$lookup": bson.M{
				"from": "components",
				"let": bson.M{
					"application_id": "$dependencies.applications",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": bson.A{"$_id", "$$application_id"},
							},
						},
					},
					{
						"$project": bson.M{
							"_id":        1.0,
							"name":       1.0,
							"folderName": 1.0,
						},
					},
				},
				"as": "dependencies.applications",
			},
		},
		{
			"$lookup": bson.M{
				"from": "functions",
				"let": bson.M{
					"function_id": "$dependencies.functions",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": bson.A{"$_id", "$$function_id"},
							},
						},
					},
					{
						"$project": bson.M{
							"_id":  1.0,
							"name": 1.0,
						},
					},
				},
				"as": "dependencies.functions",
			},
		},
		{
			"$lookup": bson.M{
				"from": "components",
				"let": bson.M{
					"universalApp_id": "$dependencies.universalApps",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": bson.A{"$_id", "$$universalApp_id"},
							},
						},
					},
					{
						"$project": bson.M{
							"_id":        1.0,
							"name":       1.0,
							"folderName": 1.0,
						},
					},
				},
				"as": "dependencies.universalApps",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$subscriber",
				"preserveNullAndEmptyArrays": false,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$subscriber.repository",
				"preserveNullAndEmptyArrays": false,
			},
		},
		{
			"$project": bson.M{
				"repository":            "$repository.name",
				"logoUrl":               "$logoUrl",
				"name":                  "$name",
				"gatewayProjectVersion": "$gateway.projectVersion",
				"isActive":              "$isActive",
				"description":           "$description",
				"autoTracking":          "$autoTracking",
				"subscriberDatabase":    "$subscriber.database.name",
				"subscriberRepository":  "$subscriber.repository.name",
				"subscriberName":        "$subscriber.fullname",
				"dependencies":          1.0,
			},
		},
	}

	cursor, err := p.getCollection().Aggregate(context.Background(), pipeline)

	if err != nil {
		logger.Write(logger.ERROR, fmt.Sprintf("error when creating cursor over result :: %s", err.Error()), logger.STDOUT)
		return "", grpc.Errorf(codes.Internal, "database error :: %s", err.Error())
	}

	if err = cursor.All(context.Background(), &projects); err != nil {
		logger.Write(logger.ERROR, fmt.Sprintf("error when iterating over result :: %s", err.Error()), logger.STDOUT)
		return "", grpc.Errorf(codes.Internal, "database error :: %s", err.Error())
	}

	if len(projects) > 0 {
		p = &projects[0]
	}

	json, err := json.Marshal(p)

	if err != nil {
		logger.Write(logger.ERROR, fmt.Sprintf("error when marshall struct into json :: %s", err.Error()), logger.STDOUT)
		return "", grpc.Errorf(codes.Internal, "marshall struct error :: %s", err.Error())
	}

	return string(json), nil
}

func (p *ProjectMiddleware) getCollection() *mongo.Collection {
	return db.MongoClient.Collection("projects")
}
