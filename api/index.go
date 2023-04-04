package api

import (
	"net/http"
	"strconv"

	"github.com/chokey2nv/obiex.finance/config"
	"github.com/gin-gonic/gin"
)

type APIClient struct {
	*config.Config
}

func NewAPIClient(cfg *config.Config) *APIClient {
	return &APIClient{
		cfg,
	}
}
func (apiClient *APIClient) GetAllTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	count, err := apiClient.DBClient.GetTransactionCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction count"})
		return
	}
	// Retrieve transactions with pagination
	txs, err := apiClient.DBClient.GetAllTransactions(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": txs,
		"meta": gin.H{
			"total_count": count,
			"page":        page,
			"limit":       limit,
		},
	})
}
func (apiClient *APIClient) GetAllTransactionsByAddress(c *gin.Context) {
	address := c.Query("address")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	count, err := apiClient.DBClient.GetTransactionCountByWallet(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction count"})
		return
	}
	// Retrieve transactions with pagination
	txs, err := apiClient.DBClient.GetAllTransactionsByAddress(address, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": txs,
		"meta": gin.H{
			"total_count": count,
			"page":        page,
			"limit":       limit,
		},
	})
}
func (apiClient *APIClient) GetAllTransactionsByClientId(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	clientId, err := strconv.ParseUint(c.Query("client_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "client id not valid"})
		return
	}

	count, err := apiClient.DBClient.GetTransactionCountByClientId(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transaction count"})
		return
	}
	// Retrieve transactions with pagination
	txs, err := apiClient.DBClient.GetAllTransactionsByClientId(clientId, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": txs,
		"meta": gin.H{
			"total_count": count,
			"page":        page,
			"limit":       limit,
		},
	})
}
