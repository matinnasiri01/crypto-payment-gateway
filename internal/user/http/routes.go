package http

import(
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Route(r *gin.RouterGroup, po *pgxpool.Pool) {

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", RegisterHandler(po))
	}


}