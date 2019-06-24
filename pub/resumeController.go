package pub

import (
	"net/http"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GetResume = func(w http.ResponseWriter, r *http.Request) {

	resumeID := r.Context().Value("resume").(primitive.ObjectID)

	res := models.GetResumeByID(resumeID)

	utils.Respond(w, res)
}
