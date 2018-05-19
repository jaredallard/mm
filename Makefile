release:
	mkdir -p "Release"
	gox -parallel=4 -ldflags="-s -w" -output="Release/{{.Dir}}-{{.OS}}-{{.Arch}}"

debug:
	go build
	mv mm ~/.bin/mm

clean:
	rm -rf Release
