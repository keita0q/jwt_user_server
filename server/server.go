package server

import (
	"github.com/pressly/chi/middleware"
	"net/http"
	"github.com/pressly/chi"
	"strings"
	"fmt"
	"sync"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
	"path/filepath"
	"context"
	"github.com/keita0q/user_server/database/applicationDatabase"
	"github.com/keita0q/user_server/auth"
	"github.com/keita0q/user_server/database/sequreDB"
	"github.com/keita0q/user_server/model"
)

type Config struct {
	ContextPath    string
	Port           int
	ClientDir      string

	Database       applicationDatabase.Database
	SequreDB       sequreDB.SequreDB
	Auth           auth.Auth
}

func Run(aConfig *Config) error {
	log.Println("User Server is Running at " + strconv.Itoa(aConfig.Port) + " Port ... ")
	tMutex := sync.RWMutex{}

	tApplicationRouter := chi.NewRouter()
	tApplicationRouter.Use(middleware.RequestID)
	tApplicationRouter.Use(middleware.RealIP)
	tApplicationRouter.Use(middleware.Logger)
	tApplicationRouter.Use(middleware.Recoverer)
	tApplicationRouter.Use(middleware.NoCache)

	tContextPath := "/" + strings.TrimSuffix(strings.TrimPrefix(aConfig.ContextPath, "/"), "/")

	tApplicationRouter.Get(tContextPath, func(aWriter http.ResponseWriter, aRequest *http.Request) {
		http.Redirect(aWriter, aRequest, tContextPath + "/", http.StatusMovedPermanently)
	})

	tApplicationRouter.Mount(tContextPath + "/", func() http.Handler {
		tMainRouter := chi.NewRouter()

		// ユーザー用UI
		tMainRouter.Get("/", func(aWriter http.ResponseWriter, aRequest *http.Request) {
			http.ServeFile(aWriter, aRequest, filepath.Join(aConfig.ClientDir, "self", "index.html"))
		})
		tMainRouter.Get("/*", func(aWriter http.ResponseWriter, aRequest *http.Request) {
			http.ServeFile(aWriter, aRequest, filepath.Join(aConfig.ClientDir, "self", strings.TrimPrefix(aRequest.URL.Path, tContextPath)))
		})

		//rest api
		tMainRouter.Mount("/rest/v1", func() http.Handler {
			tRestApiRouter := chi.NewRouter()

			// キャッシュさせない
			tRestApiRouter.Use(middleware.NoCache)

			tRestApiRouter.Post("/users", func(aWriter http.ResponseWriter, aRequest *http.Request) {
				tMutex.Lock()
				defer tMutex.Unlock()
				defer func() {
					io.Copy(ioutil.Discard, aRequest.Body)
					aRequest.Body.Close()
				}()

				tBytes, tError := ioutil.ReadAll(aRequest.Body)
				if tError != nil {
					http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				tUser := &model.User{}
				if tError := json.Unmarshal(tBytes, tUser); tError != nil {
					http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				if tError := aConfig.SequreDB.SaveUser(tUser); tError != nil {
					http.Error(aWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				aWriter.WriteHeader(http.StatusNoContent)
			})

			// login
			tRestApiRouter.Post("/login", func(aWriter http.ResponseWriter, aRequest *http.Request) {
				tMutex.Lock()
				defer tMutex.Unlock()
				defer func() {
					io.Copy(ioutil.Discard, aRequest.Body)
					aRequest.Body.Close()
				}()

				type user struct {
					ID       string `json:"id"`
					password string `json:"password"`
				}

				tBytes, tError := ioutil.ReadAll(aRequest.Body)
				if tError != nil {
					http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				tUser := &user{}
				if tError := json.Unmarshal(tBytes, tUser); tError != nil {
					http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				tToken, tError := aConfig.Auth.CreateToken(tUser.ID, tUser.password)
				if tError != nil {
					switch tError.(type){
					case *auth.NotFoundError:
						http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					default:
						http.Error(aWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
				}

				type response struct {
					Token string `json:"token"`
				}

				serveJSON(aWriter, &response{Token:tToken})
			})

			//子供用
			tRestApiRouter.Mount("/children", func() http.Handler {
				tChildRouter := chi.NewRouter()

				tChildRouter.Mount("/me", func() http.Handler {
					tMeRouter := chi.NewRouter()

					tMeRouter.Use(func(aNext http.Handler) http.Handler {
						return http.HandlerFunc(func(aWriter http.ResponseWriter, aRequest *http.Request) {
							tTokenString := aRequest.Header.Get("Authorization")

							tClaim, tOK, tError := aConfig.Auth.Authenticate(tTokenString)
							if !tOK || tError != nil {
								fmt.Println(tError)
								http.Error(aWriter, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
								return
							}

							tContext := context.WithValue(aRequest.Context(), "userID", tClaim.GetUserID())
							aNext.ServeHTTP(aWriter, aRequest.WithContext(tContext))
						})
					})

					tMeRouter.Get("/", func(aWriter http.ResponseWriter, aRequest *http.Request) {
						tMutex.Lock()
						defer tMutex.Unlock()

						defer func() {
							io.Copy(ioutil.Discard, aRequest.Body)
							aRequest.Body.Close()
						}()

						tID := aRequest.Context().Value("userID").(string)

						tUser, tError := aConfig.SequreDB.LoadUser(tID)
						if tError != nil {
							http.Error(aWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
							return
						}
						serveJSON(aWriter, tUser)
					})
					return tMeRouter
				}())

				tChildRouter.Get("/:childID", func(aWriter http.ResponseWriter, aRequest *http.Request) {
					tMutex.Lock()
					defer tMutex.Unlock()

					defer func() {
						io.Copy(ioutil.Discard, aRequest.Body)
						aRequest.Body.Close()
					}()

					tChildID := chi.URLParam(aRequest, "childID")
					tChild, tError := aConfig.Database.LoadChild(tChildID)
					if tError != nil {
						http.Error(aWriter, tError.Error(), http.StatusInternalServerError)
						return
					}
					serveJSON(aWriter, tChild)
				})

				tChildRouter.Post("/", func(aWriter http.ResponseWriter, aRequest *http.Request) {
					tMutex.Lock()
					defer tMutex.Unlock()

					defer func() {
						io.Copy(ioutil.Discard, aRequest.Body)
						aRequest.Body.Close()
					}()

					tBytes, tError := ioutil.ReadAll(aRequest.Body)
					if tError != nil {
						http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					tChild := &model.Child{}
					if tError := json.Unmarshal(tBytes, tChild); tError != nil {
						http.Error(aWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}

					if tError := aConfig.Database.SaveChild(tChild); tError != nil {
						http.Error(aWriter, tError.Error(), http.StatusInternalServerError)
						return
					}

					aWriter.WriteHeader(http.StatusNoContent)
				})

				return tChildRouter
			}())

			return tRestApiRouter
		}())
		return tMainRouter
	}())
	return http.ListenAndServe(fmt.Sprintf(":%d", aConfig.Port), tApplicationRouter)
}

func serveJSON(aWriter http.ResponseWriter, aAny interface{}) {
	tBytes, tError := json.MarshalIndent(aAny, "", "  ")
	if tError != nil {
		http.Error(aWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	aWriter.WriteHeader(http.StatusOK)
	aWriter.Header().Set("Content-Type", "application/json")
	aWriter.Write(tBytes)
}