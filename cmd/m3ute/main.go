package main

import (
	"bufio"
	"errors"
	"fmt"
	garg "github.com/alexflint/go-arg"
	//m3u "github.com/k3a/go-m3u"
	"log"
	"os"
)

type args struct {
	CacheFile          string   `arg:"-c,--cache-file" help:"Location of item JSON cache file" default:"cache_item.db"`
	CacheLoad          bool     `arg:"-C,--cache" help:"Run query to load cache; Does not produce any m3u output"`
	InputPlaylistFiles []string `arg:"positional"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := new(args)

	garg.MustParse(args)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tracks, err := collectPlaylists(args.InputPlaylistFiles)
	if err != nil {
		log.Fatal(err)
	}
	concatenate(tracks)
}

type Track struct {
	TitleEtc  string
	FileOrUrl string
}

func concatenate(tracks []*Track) {
	fmt.Println("#EXTM3U")
	for i := 0; i < len(tracks); i++ {
		fmt.Println(tracks[i].TitleEtc)
		fmt.Println(tracks[i].FileOrUrl)
	}

}

func collectPlaylists(fs []string) ([]*Track, error) {

	tracks := make([]*Track, 0)

	for i := 0; i < len(fs); i++ {
		if !fileExists(fs[i]) {
			return nil, errors.New("File:" + fs[i] + " does not exist")
		}

		file, err := os.Open(fs[i])
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		//first := true
		for scanner.Scan() {
			line1 := scanner.Text()
			if line1[0] == '#' {
				if len(line1) < 8 || line1[0:8] != "#EXTINF:" {
					continue
				}
			}
			// if first {
			// 	first = false
			// 	if !scanner.Scan() {
			// 		return nil, errors.New("Error reading file: " + fs[i])
			// 	}

			// 	line1 = scanner.Text()
			// }
			if !scanner.Scan() {
				return nil, errors.New("Error reading file: " + fs[i])
			}
			line2 := scanner.Text()
			tracks = append(tracks, &Track{TitleEtc: line1, FileOrUrl: line2})

		}
	}
	return tracks, nil
}
