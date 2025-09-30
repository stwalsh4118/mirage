package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stwalsh4118/mirageapi/internal/railway"
)

// ServicesController handles Railway service provisioning endpoints.
type ServicesController struct {
	Railway *railway.Client
}

// RegisterRoutes registers service-related routes under the provided router group.
func (c *ServicesController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/provision/services", c.ProvisionServices)
}

// ProvisionServicesRequest creates one or more services in a given environment.
type ProvisionServicesRequest struct {
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Services      []struct {
		Name   string  `json:"name"`
		Repo   *string `json:"repo"`
		Branch *string `json:"branch"`
	} `json:"services"`
	RequestID string `json:"requestId"`
}

type ProvisionServicesResponse struct {
	ServiceIDs []string `json:"serviceIds"`
}

// ProvisionServices creates services sequentially and returns their IDs.
func (c *ServicesController) ProvisionServices(ctx *gin.Context) {
	if c.Railway == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "railway client not configured"})
		return
	}
	var req ProvisionServicesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ids := make([]string, 0, len(req.Services))
	for _, s := range req.Services {
		out, err := c.Railway.CreateService(ctx, railway.CreateServiceInput{ProjectID: req.ProjectID, EnvironmentID: req.EnvironmentID, Name: s.Name, Repo: s.Repo, Branch: s.Branch})
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "service": s.Name, "partial": ids})
			return
		}
		ids = append(ids, out.ServiceID)
	}
	ctx.JSON(http.StatusOK, ProvisionServicesResponse{ServiceIDs: ids})
}




