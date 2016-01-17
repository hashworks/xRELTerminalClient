package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
	"github.com/hashworks/xRELTerminalClient/api"
)

func selectFavList(selectPrefix string) (int, error) {
	var id	int
	var err	error

	if selectPrefix == "" {
		selectPrefix = "Please choose one: "
	}

	favLists, err := api.Favs_GetLists()
	if err == nil {
		favListCount := len(favLists)
		if favListCount == 0 {
			err = errors.New("You have no favorites lists.")
		} else if favListCount > 1 {
			for i := 0; i < favListCount; i++ {
				favList := favLists[i]
				if favList.Description != "" {
					fmt.Printf("(%d) %s (%s)\n", i + 1, favList.Name, favList.Description)
				} else {
					fmt.Printf("(%d) %s\n", i + 1, favList.Name)
				}
			}
			selection := -1
			fmt.Print(selectPrefix)
			for (selection < 1 || selection > favListCount) {
				fmt.Scanln(&selection)
			}
			fmt.Println()
			id = favLists[selection - 1].Id
		} else {
			id = favLists[0].Id
		}
	}
	return id, err
}

func addEntryToFavList(extInfoId string) {
	id, err := selectFavList("Please choose the list you want to add an entry to: ")
	OK(err, "Failed to get your favorites lists:\n")
	result, err := api.Favs_ListAddEntry(strconv.Itoa(id), extInfoId)
	OK(err, "Failed to add entry:\n")
	if result.Success == 1 {
		fmt.Println("Sucessfully added \"" + result.ExtInfo.Title + "\".")
	} else {
		fmt.Println("Failed to add \"" + result.ExtInfo.Title + "\". No idea why.")
	}
}

func RemoveFavEntry() {
	id, err := selectFavList("Please choose the list you want to remove an entry from: ")
	OK(err, "Failed to get your favorites lists:\n")
	favListEntries, err := api.Favs_GetListEntries(strconv.Itoa(id), false)
	OK(err, "Failed to get favorites list entries:\n")
	favListEntriesCount := len(favListEntries)
	if favListEntriesCount == 0 {
		fmt.Println("You have no favorites list entries on this list.")
		os.Exit(1)
	} else {
		var entryId string
		if favListEntriesCount == 1 {
			var selection string
			fmt.Print("Do you really want to remove \"" + favListEntries[0].Title + "\"? (y/N) ")
			fmt.Scanln(&selection)
			if selection == "y" {
				entryId = favListEntries[0].Id
			} else {
				os.Exit(0)
			}
		} else {
			for i := 0; i < favListEntriesCount; i++ {
				favListEntry := favListEntries[i]
				fmt.Printf("(%d) [%s] %s\n", i + 1, strings.ToUpper(favListEntry.Type), favListEntry.Title)
			}
			selection := -1
			fmt.Print("Please choose the entry you want to remove: ")
			for (selection < 1 || selection > favListEntriesCount) {
				fmt.Scanln(&selection)
			}
			entryId = favListEntries[selection-1].Id
		}
		fmt.Println()
		result, err := api.Favs_ListRemoveEntry(strconv.Itoa(id), entryId)
		OK(err, "Failed to remove entry:\n")
		if result.Success == 1 {
			fmt.Println("Sucessfully removed \"" + result.ExtInfo.Title + "\".")
		} else {
			fmt.Println("Failed to remove \"" + result.ExtInfo.Title + "\". No idea why.")
		}
	}
}

func ShowFavEntries() {
	id, err := selectFavList("")
	OK(err, "Failed to get your favorites lists:\n")
	favListEntries, err := api.Favs_GetListEntries(strconv.Itoa(id), true)
	OK(err, "Failed to get favorites list entries:\n")
	favListEntriesCount := len(favListEntries)
	if favListEntriesCount == 0 {
		fmt.Println("You have no favorites list entries on this list.")
		os.Exit(1)
	} else {
		for i := 0; i < favListEntriesCount; i++ {
			entry := favListEntries[i]
			fmt.Println(entry.Title + " [" + strings.ToUpper(entry.Type) + "]")
			releaseCount := len(entry.Releases)
			p2pReleaseCount := len(entry.P2PReleases)
			if releaseCount == 0 && p2pReleaseCount == 0 {
				fmt.Println("\tNo new releases.")
			} else {
				for i := 0; i < releaseCount; i++ {
					release := entry.Releases[i]
					fmt.Println("\t[SCENE] " + release.Dirname + " (" + release.LinkHref + ")")
				}
				for i := 0; i < p2pReleaseCount; i++ {
					release := entry.P2PReleases[i]
					fmt.Println("\t[P2P]   " + release.Dirname + " (" + release.LinkHref + ")")
				}
			}
			fmt.Println()
		}
	}
}