package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ViniNepo/church-player-backend/internal/config"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/fileoperator"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/handler"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
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
	//executableDir := filepath.Dir(executable)
	executableDir := "C:\\Users\\vinep\\Downloads\\ChurchPlayer\\ChurchPlayer\\resources\\app\\"

	// Constrói o caminho para o arquivo que está na mesma pasta do binário
	dbPath := filepath.Join(executableDir, "church.db")
	log.Println(executable)
	log.Println(dbPath)

	// Configurar o roteador Gorilla Mux
	router := mux.NewRouter()

	fileManager := fileoperator.NewFileManager(executableDir)

	//rotas
	router.HandleFunc("/ping", pong).Methods("GET")

	router.HandleFunc("/file-control/{path}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		fileControlUseCase := usecase.NewFileControlUseCase(fileManager)
		fileControl := handler.NewFileControlHandler(fileControlUseCase)
		fileControl.Play(w, r)
	}).Methods("GET")

	router.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
		albumHandler := handler.NewAlbumHandler(albumUseCase)

		albumHandler.Create(w, r)
	}).Methods("POST")

	router.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
		albumHandler := handler.NewAlbumHandler(albumUseCase)

		albumHandler.FindAll(w, r)
	}).Methods("GET")

	router.HandleFunc("/albums/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
		albumHandler := handler.NewAlbumHandler(albumUseCase)

		albumHandler.FindByID(w, r)
	}).Methods("GET")

	router.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
		albumHandler := handler.NewAlbumHandler(albumUseCase)

		albumHandler.Update(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/albums/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		albumUseCase := usecase.NewAlbumUseCase(albumDB, songDB, fileManager)
		albumHandler := handler.NewAlbumHandler(albumUseCase)

		albumHandler.Delete(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/songs/{album_id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.Create(w, r)
	}).Methods("POST")

	router.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.FindAll(w, r)
	}).Methods("GET")

	router.HandleFunc("/songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.FindById(w, r)
	}).Methods("GET")

	router.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.Update(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/songs/orders", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.UpdateOrder(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		songDB := sqlite.NewSongDB(db)
		songUseCase := usecase.NewSongUseCase(songDB, fileManager)
		songHandler := handler.NewSongHandler(songUseCase)

		songHandler.Delete(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/moments", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.Create(w, r)
	}).Methods("POST")

	router.HandleFunc("/moments/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.FindByID(w, r)
	}).Methods("GET")

	router.HandleFunc("/moments/{id}/all", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.FindAllBySectionID(w, r)
	}).Methods("GET")

	router.HandleFunc("/moments", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.Update(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/moments/orders", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.UpdateOrder(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/moments/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		momentDB := sqlite.NewMomentDB(db)
		momentUseCase := usecase.NewMomentUseCase(momentDB)
		momentHandler := handler.NewMomentHandler(momentUseCase)

		momentHandler.Delete(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/sections", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		sectionDB := sqlite.NewSectionDB(db)
		sectionUseCase := usecase.NewSectionUseCase(sectionDB)
		sectionHandler := handler.NewSectionHandler(sectionUseCase)

		sectionHandler.Create(w, r)
	}).Methods("POST")

	router.HandleFunc("/sections/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		sectionDB := sqlite.NewSectionDB(db)
		sectionUseCase := usecase.NewSectionUseCase(sectionDB)
		sectionHandler := handler.NewSectionHandler(sectionUseCase)

		sectionHandler.FindAllByWorshipID(w, r)
	}).Methods("GET")

	router.HandleFunc("/sections", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		sectionDB := sqlite.NewSectionDB(db)
		sectionUseCase := usecase.NewSectionUseCase(sectionDB)
		sectionHandler := handler.NewSectionHandler(sectionUseCase)

		sectionHandler.Update(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/sections/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		sectionDB := sqlite.NewSectionDB(db)
		sectionUseCase := usecase.NewSectionUseCase(sectionDB)
		sectionHandler := handler.NewSectionHandler(sectionUseCase)

		sectionHandler.Delete(w, r)
	}).Methods("DELETE")

	router.HandleFunc("/worships", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		momentDB := sqlite.NewMomentDB(db)
		sectionDB := sqlite.NewSectionDB(db)
		worshipDB := sqlite.NewWorshipDB(db)
		worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
		worshipHandler := handler.NewWorshipHandler(worshipUseCase)
		worshipHandler.Create(w, r)
	}).Methods("POST")

	router.HandleFunc("/worships", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		momentDB := sqlite.NewMomentDB(db)
		sectionDB := sqlite.NewSectionDB(db)
		worshipDB := sqlite.NewWorshipDB(db)
		worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
		worshipHandler := handler.NewWorshipHandler(worshipUseCase)
		worshipHandler.FindAll(w, r)
	}).Methods("GET")

	router.HandleFunc("/worships/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		momentDB := sqlite.NewMomentDB(db)
		sectionDB := sqlite.NewSectionDB(db)
		worshipDB := sqlite.NewWorshipDB(db)
		worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
		worshipHandler := handler.NewWorshipHandler(worshipUseCase)
		worshipHandler.FindByID(w, r)
	}).Methods("GET")

	router.HandleFunc("/worships", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		momentDB := sqlite.NewMomentDB(db)
		sectionDB := sqlite.NewSectionDB(db)
		worshipDB := sqlite.NewWorshipDB(db)
		worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
		worshipHandler := handler.NewWorshipHandler(worshipUseCase)
		worshipHandler.Update(w, r)
	}).Methods("PATCH")

	router.HandleFunc("/worships/{id}", func(w http.ResponseWriter, r *http.Request) {
		db, err := openDB(dbPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		albumDB := sqlite.NewAlbumDB(db)
		songDB := sqlite.NewSongDB(db)
		momentDB := sqlite.NewMomentDB(db)
		sectionDB := sqlite.NewSectionDB(db)
		worshipDB := sqlite.NewWorshipDB(db)
		worshipUseCase := usecase.NewWorshipUseCase(worshipDB, momentDB, songDB, albumDB, sectionDB, fileManager)
		worshipHandler := handler.NewWorshipHandler(worshipUseCase)
		worshipHandler.Delete(w, r)
	}).Methods("DELETE")

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

func openDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco de dados: %w", err)
	}
	return db, nil
}
