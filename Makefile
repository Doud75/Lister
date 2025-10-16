-include .env
export

.DEFAULT_GOAL := help

include make/*.mk

help:
	@echo "Usage: make [cible]"
	@echo ""
	@echo "Cibles principales:"
	@echo "  up             Démarre tous les services Docker."
	@echo "  down           Arrête tous les services Docker."
	@echo "  logs           Affiche les logs des services."
	@echo "  test-all       Lance tous les tests (unitaires et E2E)."
	@echo "  test           Lance uniquement les tests E2E."
	@echo "  test-unit      Lance uniquement les tests unitaires."
	@echo ""
	@echo "Consultez les fichiers dans le dossier 'make/' pour voir toutes les cibles disponibles."

.PHONY: help