.PHONY: run

run:
	cd backend && go run main.go & \
	cd frontend/flutter_app && flutter run -d chrome
# command = make run