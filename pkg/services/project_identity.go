package services

import (
	"context"
	"encoding/json"
	"fmt"
	"lucas-stellet/api-grpc-mongodb/pkg/db"
	"lucas-stellet/api-grpc-mongodb/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Tenant struct {
	ID           string `bson:"_id"`
	ClientID     string `bson:"client_id"`
	ClientSecret string `bson:"client_secret"`
}

type credential map[string]string

type Credentials []credential

type Tenants []Tenant

type ProjectIdentity struct {
	IsActive    bool        `json:"isActive" bson:"isActive"`
	Repository  string      `json:"repository" bson:"repository"`
	Subscriber  string      `json:"subscriber" bson:"subscriber"`
	Project     string      `json:"project" bson:"project"`
	Auth        string      `json:"auth" bson:"auth"`
	Tenants     Tenants     `json:",omitempty" bson:"tenants"`
	Credentials Credentials `json:"credentials"`
}

func (p *ProjectIdentity) GetData(gatewayName string, apiVersion int32) (string, error) {
	var projects []ProjectIdentity

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"url":                gatewayName,
				"deleted":            false,
				"type":               "EXPOSE_API",
				"gateway.apiVersion": apiVersion,
			},
		},
		{
			"$lookup": bson.M{
				"from": "tenants",
				"let": bson.M{
					"tenant_id": "$tenants",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": bson.A{"$_id", "$$tenant_id"},
							},
						},
					},
					{
						"$project": bson.M{
							"client_id":     1.0,
							"client_secret": 1.0,
						},
					},
				},
				"as": "tenants",
			},
		},
		{
			"$project": bson.M{
				"isActive":   "$isActive",
				"repository": "$repository.name",
				"subscriber": "$subscriber",
				"project":    "$_id",
				"auth":       "$gateway.authentication",
				"tenants":    1,
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

	p.generateCredentials()

	json, err := json.Marshal(p)

	if err != nil {
		logger.Write(logger.ERROR, fmt.Sprintf("error when marshall struct into json :: %s", err.Error()), logger.STDOUT)
		return "", grpc.Errorf(codes.Internal, "marshall struct error :: %s", err.Error())
	}

	return string(json), nil
}

func (p *ProjectIdentity) getCollection() *mongo.Collection {
	return db.MongoClient.Collection("projects")
}

func (p *ProjectIdentity) generateCredentials() {
	for _, t := range p.Tenants {
		p.Credentials = append(p.Credentials, p.generateCredentialFromTenant(t.ID, t.ClientID, t.ClientSecret))
	}

	p.removeTenants()
}

func (p *ProjectIdentity) removeTenants() {
	p.Tenants = p.Tenants[:0]
}

func (p *ProjectIdentity) generateCredentialFromTenant(tenantID, clientID, clientSecret string) credential {
	var cred credential
	switch p.Auth {
	case "API_KEY":
		cred = map[string]string{
			fmt.Sprintf("cred:%s", clientID): tenantID,
		}
	default:
		cred = map[string]string{
			fmt.Sprintf("cred:%s:%s", clientID, clientSecret): tenantID,
		}
	}
	return cred
}
