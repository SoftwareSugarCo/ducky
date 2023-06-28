package util

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
)

var duckyBig = []byte(`
  _____  _    _  _____ _  ____     __
 |  __ \| |  | |/ ____| |/ /\ \   / /
 | |  | | |  | | |    | ' /  \ \_/ / 
 | |  | | |  | | |    |  <    \   /  
 | |__| | |__| | |____| . \    | |   
 |_____/ \____/ \_____|_|\_\   |_|   
`)

var duckyBox = []byte(` 

.----------------.  .----------------.  .----------------.  .----------------.  .----------------. 
| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |
| |  ________    | || | _____  _____ | || |     ______   | || |  ___  ____   | || |  ____  ____  | |
| | |_   ___ '.  | || ||_   _||_   _|| || |   .' ___  |  | || | |_  ||_  _|  | || | |_  _||_  _| | |
| |   | |   '. \ | || |  | |    | |  | || |  / .'   \_|  | || |   | |_/ /    | || |   \ \  / /   | |
| |   | |    | | | || |  | '    ' |  | || |  | |         | || |   |  __'.    | || |    \ \/ /    | |
| |  _| |___.' / | || |   \ '--' /   | || |  \ '.___.'\  | || |  _| |  \ \_  | || |    _|  |_    | |
| | |________.'  | || |    '.__.'    | || |   '._____.'  | || | |____||____| | || |   |______|   | |
| |              | || |              | || |              | || |              | || |              | |
| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |
 '----------------'  '----------------'  '----------------'  '----------------'  '----------------' 
`)

var duckyLetters = []byte(`
DDDDDDDDDDDDD       UUUUUUUU     UUUUUUUU       CCCCCCCCCCCCCKKKKKKKKK    KKKKKKKYYYYYYY       YYYYYYY
D::::::::::::DDD    U::::::U     U::::::U    CCC::::::::::::CK:::::::K    K:::::KY:::::Y       Y:::::Y
D:::::::::::::::DD  U::::::U     U::::::U  CC:::::::::::::::CK:::::::K    K:::::KY:::::Y       Y:::::Y
DDD:::::DDDDD:::::D UU:::::U     U:::::UU C:::::CCCCCCCC::::CK:::::::K   K::::::KY::::::Y     Y::::::Y
  D:::::D    D:::::D U:::::U     U:::::U C:::::C       CCCCCCKK::::::K  K:::::KKKYYY:::::Y   Y:::::YYY
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K:::::K K:::::K      Y:::::Y Y:::::Y   
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K::::::K:::::K        Y:::::Y:::::Y    
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K:::::::::::K          Y:::::::::Y     
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K:::::::::::K           Y:::::::Y      
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K::::::K:::::K           Y:::::Y       
  D:::::D     D:::::DU:::::D     D:::::UC:::::C                K:::::K K:::::K          Y:::::Y       
  D:::::D    D:::::D U::::::U   U::::::U C:::::C       CCCCCCKK::::::K  K:::::KKK       Y:::::Y       
DDD:::::DDDDD:::::D  U:::::::UUU:::::::U  C:::::CCCCCCCC::::CK:::::::K   K::::::K       Y:::::Y       
D:::::::::::::::DD    UU:::::::::::::UU    CC:::::::::::::::CK:::::::K    K:::::K    YYYY:::::YYYY    
D::::::::::::DDD        UU:::::::::UU        CCC::::::::::::CK:::::::K    K:::::K    Y:::::::::::Y    
DDDDDDDDDDDDD             UUUUUUUUU             CCCCCCCCCCCCCKKKKKKKKK    KKKKKKK    YYYYYYYYYYYYY
`)

var duckyRec = []byte(`
 ____  _____ _____ _____ __ __ 
|    \|  |  |     |  |  |  |  |
|  |  |  |  |   --|    -|_   _|
|____/|_____|_____|__|__| |_|
`)

var duckyBlood = []byte(`
▓█████▄  █    ██  ▄████▄   ██ ▄█▀▓██   ██▓
▒██▀ ██▌ ██  ▓██▒▒██▀ ▀█   ██▄█▒  ▒██  ██▒
░██   █▌▓██  ▒██░▒▓█    ▄ ▓███▄░   ▒██ ██░
░▓█▄   ▌▓▓█  ░██░▒▓▓▄ ▄██▒▓██ █▄   ░ ▐██▓░
░▒████▓ ▒▒█████▓ ▒ ▓███▀ ░▒██▒ █▄  ░ ██▒▓░
 ▒▒▓  ▒ ░▒▓▒ ▒ ▒ ░ ░▒ ▒  ░▒ ▒▒ ▓▒   ██▒▒▒ 
 ░ ▒  ▒ ░░▒░ ░ ░   ░  ▒   ░ ░▒ ▒░ ▓██ ░▒░ 
 ░ ░  ░  ░░░ ░ ░ ░        ░ ░░ ░  ▒ ▒ ░░  
   ░       ░     ░ ░      ░  ░    ░ ░     
 ░               ░                ░ ░
`)

var duckyCalvin = []byte(`
╔╦╗╦ ╦╔═╗╦╔═╦ ╦
 ║║║ ║║  ╠╩╗╚╦╝
═╩╝╚═╝╚═╝╩ ╩ ╩
`)

var duckyCorps = []byte(`
████████▄  ███    █▄   ▄████████    ▄█   ▄█▄ ▄██   ▄   
███   ▀███ ███    ███ ███    ███   ███ ▄███▀ ███   ██▄ 
███    ███ ███    ███ ███    █▀    ███▐██▀   ███▄▄▄███ 
███    ███ ███    ███ ███         ▄█████▀    ▀▀▀▀▀▀███ 
███    ███ ███    ███ ███        ▀▀█████▄    ▄██   ███ 
███    ███ ███    ███ ███    █▄    ███▐██▄   ███   ███ 
███   ▄███ ███    ███ ███    ███   ███ ▀███▄ ███   ███ 
████████▀  ████████▀  ████████▀    ███   ▀█▀  ▀█████▀
`)

var ducky3D = []byte(`
 *******   **     **   ******  **   ** **    **
/**////** /**    /**  **////**/**  ** //**  ** 
/**    /**/**    /** **    // /** **   //****  
/**    /**/**    /**/**       /****     //**   
/**    /**/**    /**/**       /**/**     /**   
/**    ** /**    /**//**    **/**//**    /**   
/*******  //*******  //****** /** //**   /**   
///////    ///////    //////  //   //    //
`)

var duckyAMC = []byte(`

                         ______                           
|'''''''. |         |  .~      ~. |    ..'' ''..     ..'' 
|       | |         | |           |..''         ''.''     
|       | |         | |           |''..           |       
|......'  '._______.'  '.______.' |    ''..       |
`)

var duckyUSA = []byte(`
 :::====  :::  === :::===== :::  === ::: ===
 :::  === :::  === :::      ::: ===  ::: ===
 ===  === ===  === ===      ======    ===== 
 ===  === ===  === ===      === ===    ===  
 =======   ======   ======= ===  ===   ===  
`)

func PrintDuckyHeader() {
	yellowStr := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	duckyHeaders := [][]byte{duckyBig, duckyBox, duckyLetters, duckyRec, duckyBlood, duckyCalvin, duckyCorps, ducky3D, duckyAMC, duckyUSA}

	selectedHeader := duckyHeaders[rand.Intn(len(duckyHeaders))]

	fmt.Println(yellowStr(string(selectedHeader)))

	//fmt.Println(string(selectedHeader))
}

func PrintBox(label string, items map[string]string) {
	var maxLength int

	// Find the length of the longest setting key or value
	for key, value := range items {
		if len(key) > maxLength {
			maxLength = len(key)
		}

		if len(value) > maxLength {
			maxLength = len(value)
		}
	}

	// Add two to the maxLength to account for the padding on either side of the value
	maxLength += 2

	// Print the box
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Print the label line
	fmt.Printf("| %s%-*s |\n", label, ((2*maxLength)-len(label))+3, "")

	// Print the separator line
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Print the items
	for key, value := range items {
		fmt.Printf("| %-*s | %-*s |\n", maxLength, key, maxLength, value)
	}

	// Print the bottom of the box
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	fmt.Printf("\n")
}
