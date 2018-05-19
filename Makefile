release:
	mkdir -p "Release"
	gox -parallel=4 -ldflags="-s -w" -output="Release/{{.Dir}}-{{.OS}}-{{.Arch}}"

debug:
	go build
	mv hpkg ~/.bin/hpkg

clean:
	rm -rf Release