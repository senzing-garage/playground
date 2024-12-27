package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataSourceKey struct {
	Data_Source string
}

type Record struct {
	Data_Source string
	Record_ID   string
}

var (
	ctx               = context.TODO()
	grpcAddress       = "localhost:8261"
	jsonDataSource    DataSourceKey
	homePath          = "./"
	jsonRecord        Record
	truthSetURLPrefix = "https://raw.githubusercontent.com/Senzing/truth-sets/refs/heads/main/truthsets/demo/"
	truthSetFileNames = []string{"customers.json", "reference.json", "watchlist.json"}
)

func testErr(err error) {
	if err != nil {
		panic(err)
	}
}

func downloadFile(url string, filepath string) error {
	outputFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func extractDataSources(filePath string) []string {
	result := []string{}
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &jsonDataSource)
		testErr(err)
		if !slices.Contains(result, jsonDataSource.Data_Source) {
			result = append(result, jsonDataSource.Data_Source)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func addDatasourcesToSenzingConfig(szAbstractFactory senzing.SzAbstractFactory, dataSourceNames []string) error {
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		return err
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		return err
	}

	oldConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		return err
	}

	oldJsonConfig, err := szConfigManager.GetConfig(ctx, oldConfigID)
	if err != nil {
		return err
	}

	configHandle, err := szConfig.ImportConfig(ctx, oldJsonConfig)
	if err != nil {
		return err
	}

	for _, value := range dataSourceNames {
		_, err := szConfig.AddDataSource(ctx, configHandle, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	newJsonConfig, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return err
	}

	newConfigID, err := szConfigManager.AddConfig(ctx, newJsonConfig, "Add TruthSet datasources")
	if err != nil {
		return err
	}

	err = szConfigManager.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	if err != nil {
		return err
	}

	err = szAbstractFactory.Reinitialize(ctx, newConfigID)
	if err != nil {
		return err
	}

	return nil
}

func addRecords(szAbstractFactory senzing.SzAbstractFactory, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &jsonRecord)
		if err != nil {
			return err
		}
		result, err := szEngine.AddRecord(ctx, jsonRecord.Data_Source, jsonRecord.Record_ID, string(line), senzing.SzWithInfo)
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}

func asPrettyJSON(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}

func main() {

	// Download truth-sets files.

	for i := 0; i < len(truthSetFileNames); i++ {
		url := fmt.Sprintf("%s/%s", truthSetURLPrefix, truthSetFileNames[i])
		filepath := fmt.Sprintf("%s%s", homePath, truthSetFileNames[i])
		err := downloadFile(url, filepath)
		testErr(err)
	}

	// Create an abstract factory for acessing Senzing via gRPC.

	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	testErr(err)
	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}

	// Discover DATA_SOURCE values in records.

	dataSources := []string{}
	for _, value := range truthSetFileNames {
		partialDataSources := extractDataSources(fmt.Sprintf("%s%s", homePath, value))
		dataSources = append(dataSources, partialDataSources...)
	}
	fmt.Printf("Found the following DATA_SOURCE values in the data: %v\n", dataSources)

	// Update Senzing configuration.

	err = addDatasourcesToSenzingConfig(szAbstractFactory, dataSources)

	// Add records.

	for _, value := range truthSetFileNames {
		err = addRecords(szAbstractFactory, fmt.Sprintf("%s%s", homePath, value))
	}

	// Retrieve an entity by identifying a record of the entity. Use the `SZ_ENTITY_INCLUDE_RECORD_SUMMARY` flag from among the get_entity flags.

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	testErr(err)

	customer1070Entity, err := szEngine.GetEntityByRecordID(ctx, "CUSTOMERS", "1070", senzing.SzEntityIncludeRecordSummary)
	testErr(err)
	fmt.Println(asPrettyJSON(customer1070Entity))

	// Search for entities by attributes.

	searchProfile := ""
	searchQuery := `{
    "name_full": "robert smith",
    "date_of_birth": "11/12/1978"
}`

	searchResult, err := szEngine.SearchByAttributes(ctx, searchQuery, searchProfile, senzing.SzSearchByAttributesDefaultFlags)
	testErr(err)
	fmt.Println(asPrettyJSON(searchResult))

}
