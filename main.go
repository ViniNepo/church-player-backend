package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ViniNepo/church-player-backend/internal/config"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/fileoperator"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/handler"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Carregar a configuração
	cfg := config.LoadConfig()

	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
		return
	}

	// Obtém o diretório onde o binário está localizado
	executableDir := filepath.Dir(executable)

	// Constrói o caminho para o arquivo que está na mesma pasta do binário
	dbPath := filepath.Join(executableDir, "church.db")

	log.Print(executableDir)
	log.Print(dbPath)

	// Conectar ao banco de dados SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Configurar o roteador Gorilla Mux
	router := mux.NewRouter()

	// Configura o banco de dados
	albumDB := sqlite.NewAlbumDB(db)
	songDB := sqlite.NewSongDB(db)
	momentDB := sqlite.NewMomentDB(db)
	sectionDB := sqlite.NewSectionDB(db)
	worshipDB := sqlite.NewWorshipDB(db)

	fileManager := fileoperator.NewFileManager(executableDir)

	// Configurar os casos de uso
	albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
	songUseCase := usecase.NewSongUseCase(songDB, fileManager)
	momentUseCase := usecase.NewMomentUseCase(momentDB)
	sectionUseCase := usecase.NewSectionUseCase(sectionDB)
	worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
	fileControlUseCase := usecase.NewFileControlUseCase(fileManager)

	// Iniciar handlers
	albumHandler := handler.NewAlbumHandler(albumUseCase)
	songHandler := handler.NewSongHandler(songUseCase)
	momentHandler := handler.NewMomentHandler(momentUseCase)
	sectionHandler := handler.NewSectionHandler(sectionUseCase)
	worshipHandler := handler.NewWorshipHandler(worshipUseCase)
	fileControl := handler.NewFileControlHandler(fileControlUseCase)

	//rotas
	router.HandleFunc("/ping", pong).Methods("GET")

	router.HandleFunc("/file-control/{path}", fileControl.Play).Methods("GET")

	router.HandleFunc("/albums", albumHandler.Create).Methods("POST")
	router.HandleFunc("/albums", albumHandler.FindAll).Methods("GET")
	router.HandleFunc("/albums/{id}", albumHandler.FindByID).Methods("GET")
	router.HandleFunc("/albums", albumHandler.Update).Methods("PATCH")
	router.HandleFunc("/albums/{id}", albumHandler.Delete).Methods("DELETE")

	router.HandleFunc("/songs/{album_id}", songHandler.Create).Methods("POST")
	router.HandleFunc("/songs", songHandler.FindAll).Methods("GET")
	router.HandleFunc("/songs/{id}", songHandler.FindById).Methods("GET")
	router.HandleFunc("/songs", songHandler.Update).Methods("PATCH")
	router.HandleFunc("/songs/orders", songHandler.UpdateOrder).Methods("PATCH")
	router.HandleFunc("/songs/{id}", songHandler.Delete).Methods("DELETE")

	router.HandleFunc("/moments", momentHandler.Create).Methods("POST")
	router.HandleFunc("/moments/{id}", momentHandler.FindByID).Methods("GET")
	router.HandleFunc("/moments/{id}/all", momentHandler.FindAllBySectionID).Methods("GET")
	router.HandleFunc("/moments", momentHandler.Update).Methods("PATCH")
	router.HandleFunc("/moments/orders", momentHandler.UpdateOrder).Methods("PATCH")
	router.HandleFunc("/moments/{id}", momentHandler.Delete).Methods("DELETE")

	router.HandleFunc("/sections", sectionHandler.Create).Methods("POST")
	router.HandleFunc("/sections/{id}", sectionHandler.FindAllByWorshipID).Methods("GET")
	router.HandleFunc("/sections", sectionHandler.Update).Methods("PATCH")
	router.HandleFunc("/sections/{id}", sectionHandler.Delete).Methods("DELETE")

	router.HandleFunc("/worships", worshipHandler.Create).Methods("POST")
	router.HandleFunc("/worships", worshipHandler.FindAll).Methods("GET")
	router.HandleFunc("/worships/{id}", worshipHandler.FindByID).Methods("GET")
	router.HandleFunc("/worships", worshipHandler.Update).Methods("PATCH")
	router.HandleFunc("/worships/{id}", worshipHandler.Delete).Methods("DELETE")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Configurar Negroni como middleware
	n := negroni.Classic()
	n.UseHandler(router)

	// Iniciar o servidor
	err = http.ListenAndServe(":"+cfg.ServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(n))
	if err != nil {
		log.Println("error starting server:", err)
		return
	}
	log.Println("Server started on port", cfg.ServerPort)
}

func pong(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		return
	}
}
