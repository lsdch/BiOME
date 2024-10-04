package person

import (
	"context"
	"darco/proto/controllers"
	accounts "darco/proto/controllers/users"
	"darco/proto/db"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"
	"net/url"

	"darco/proto/models/people"
	_ "darco/proto/models/validations"

	"github.com/danielgtaylor/huma/v2"
)

type UUIDResolver struct {
	controllers.UUIDInput
	*people.Person
}

func (i *UUIDResolver) Resolve(ctx huma.Context) []error {
	person, err := people.FindPerson(db.Client(), i.ID)
	if err != nil {
		return []error{huma.Error404NotFound("Person not found")}
	}
	i.Person = &person
	return nil
}

func RegisterRoutes(r router.Router) {
	personsAPI := r.RouteGroup("/persons").
		WithTags([]string{"People", "Person"})

	router.Register(personsAPI, "ListPersons",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List persons",
			Errors:  []int{500},
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](people.ListPersons),
	)

	router.Register(personsAPI, "CreatePerson",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create person",
			Errors:  []int{400, 500},
		},
		controllers.CreateHandler[people.PersonInput, people.Person])

	router.Register(personsAPI, "UpdatePerson",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodPatch,
			Summary: "Update person",
			Errors:  []int{400, 500},
		}, controllers.UpdateByIDHandler[people.PersonUpdate](people.FindPerson))

	router.Register(personsAPI, "DeletePerson",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodDelete,
			Summary: "Delete person",
			Errors:  []int{400, 500},
		},
		controllers.DeleteByIDHandler(people.DeletePerson))

	router.Register(personsAPI, "InvitePerson",
		huma.Operation{
			Path:        "/{id}/invite",
			Method:      http.MethodPost,
			Summary:     "Invite person",
			Description: "Sends an invitation link to a person at the address provided in `dest`, allowing them to register an account assigned with a specified `role`.",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
		}, InvitePerson(accounts.InvitationClaimPath))
}

type InvitationInput struct {
	Handler *url.URL `json:"handler,omitempty" format:"uri-template" example:"http://example.com/register/{token}" doc:"A URL template with a {token} parameter, which implements the UI to validate the invitation token and fill a registration form."`
	people.InvitationOptions
}

type InvitePersonInput struct {
	resolvers.AuthRequired
	resolvers.HostResolver
	UUIDResolver
	Body InvitationInput
}

type InvitationLink struct {
	Link *url.URL `json:"invitation_link" doc:"The generated URL containing a registration token that can be shared to the invitee."`
}

type InvitePersonOutput struct {
	Body InvitationLink
}

func InvitePerson(defaultInvitationClaimPath string) router.Endpoint[InvitePersonInput, InvitePersonOutput] {
	return func(ctx context.Context, input *InvitePersonInput) (*InvitePersonOutput, error) {

		invitation, err := input.Person.
			CreateInvitation(input.Body.InvitationOptions).
			Save(input.DB())
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to register invitation token", err)
		}

		var targetURL = input.GenerateURL(defaultInvitationClaimPath)
		if input.Body.Handler != nil {
			targetURL = *input.Body.Handler
		}

		activationURL, err := invitation.Send(targetURL)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to send invitation email", err)
		}
		return &InvitePersonOutput{Body: InvitationLink{activationURL}}, nil
	}
}
