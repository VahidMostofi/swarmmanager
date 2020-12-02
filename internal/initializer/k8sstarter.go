package initializer

import(
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/viper"
)


// StartK8sAutoConfigurer ...
func StartK8sAutoConfigurer(strategy strategies.Configurer, strategyName string){
	if strings.Contains(configs.GetConfig().Results.Path, "$STRATEGY") {
		configs.GetConfig().Results.Path = strings.Replace(configs.GetConfig().Results.Path, "$STRATEGY", strategyName, 1)
		log.Println("Updating result path to", configs.GetConfig().Results.Path)
	}

	// creating directories for ResultDirectoryPath
	if err := os.MkdirAll(filepath.Dir(configs.GetConfig().Results.Path), 0777); err != nil {
		log.Panic(err)
	}

	workloadStr := viper.GetString("workloadStr")

	var lg = GetTheLoadGenerator(workloadStr)
	db := GetNewDatabase()
	a := autoconfigure.NewK8sAutoConfigurer(lg, strategy, workloadStr, db)
	log.Println("name of the test is:", viper.GetString("testName"))
	a.Start(viper.GetString("testName"), strings.Join(os.Args, " "))
	fmt.Println("done")
}