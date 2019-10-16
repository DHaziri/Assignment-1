package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "strconv"
    "time"
)

      // Global declatation of StartTime,
      //  so both main and Diagnostics can access it
var StartTime time.Time




      // Struct to get from restcountries
type CountryInfo struct {
  Code string `json:"alpha2Code"`
  Name string `json:"name"`
  Flag string `json:"flag"`
}

      // Struct to get an array of species from gbif
type CountrySpecies struct {
  CouSpe []CouSpe `json:"results"`
}
      // Each induvidual species and species key
type CouSpe struct {
  Species string `json:"species"`
  SpeciesKey int `json:"speciesKey"`
}

      // Struct to combine information from both
      //  restcountries and gbif into one
type Country struct {
  Code string `json:"alpha2Code"`
  Name string `json:"name"`
  Flag string `json:"flag"`
  Species []string `json:"species"`
  SpeciesKey []int `json:"speciesKey"`
}




      // Species struct
type Species struct {
  Key int `json:"key"`
  Kingdom string `json:"kingdom"`
  Phylum string `json:"phylum"`
  Order string `json:"order"`
  Family string `json:"family"`
  Genus string `json:"genus"`
  ScientificName string `json:"scientificName"`
  CanonicalName string `json:"canonicalName"`
  Year string `json:"year"`
}




      // Diagnostics struct
type Diagnostics struct {
   GBIF int `json:"gbif"`
   RestCountries int `json:"restcountries"`
   Version string `json:"v1"`
   UpTime int `json:"uptime"`
}




      // Gets and returns every species on the requested country
func country(w http.ResponseWriter, r *http.Request){
          // API urls
	urlCI := "http://restcountries.eu/rest/v2/alpha/"
  urlCS := "http://api.gbif.org/v1/occurrence/search?country="

          // Splits the user given url
  splits := strings.Split(r.URL.Path, "/")

          // Adds the country code to the urls,
  urlCI += splits[4]
  urlCS += strings.ToUpper(splits[4])

  limit := 20 // Matching the default limit

          // If user asks for a bigger or lower limit
  if len(splits) >= 6 && splits[5] != ""   {
          // converts the string containig the limit into an int
    limit, _ = strconv.Atoi(splits[5])
          // adds the limit to gbif url
    urlCS += "&limit="
    urlCS += splits[5]
  }


          // Gets the url for countries
	resCountryInfo, err := http.Get(urlCI)
  if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
  }
          // reads in the api for countries
	resCIData, err := ioutil.ReadAll(resCountryInfo.Body)
  if err != nil {
			log.Fatal(err)
  }
          // declaring a struct and filling inn the json and api for countries
	var resCIObject CountryInfo
	json.Unmarshal(resCIData, &resCIObject)


          // Gets the url for species
  resCountrySpecies, err := http.Get(urlCS)
	if err != nil {
			panic(err)
	}
          // reads in the api for species
	resCSData, err := ioutil.ReadAll(resCountrySpecies.Body)
	if err != nil {
			panic(err)
	}
          // declaring a struct and filling inn the json and api for species
  var resCSObject CountrySpecies
  json.Unmarshal(resCSData, &resCSObject)


          // Declaring and "transfer" the combined struct
  var responseCountry Country
  responseCountry.Code = resCIObject.Code
  responseCountry.Name = resCIObject.Name
  responseCountry.Flag = resCIObject.Flag
  for i := 0; i < limit && i < len(resCSObject.CouSpe); i++ {
    responseCountry.Species = append(responseCountry.Species, resCSObject.CouSpe[i].Species)
    responseCountry.SpeciesKey = append(responseCountry.SpeciesKey, resCSObject.CouSpe[i].SpeciesKey)
  }


          // Declaring and converting the struct into bytes
  countryBytes, err := json.Marshal(responseCountry)
          // writing out the bytes and status
  w.Header().Set("Content-Type", "application/json")
  w.Write(countryBytes)
  w.WriteHeader(http.StatusOK)
}




      // Gets and returns data on requested specie
func species(w http.ResponseWriter, r *http.Request){
          // API urls
  url := "http://api.gbif.org/v1/species/"
  urlYear := "http://api.gbif.org/v1/species/"

          // Splits the user given url
  splits := strings.Split(r.URL.Path, "/")

          // Adds the specieskey to the urls,
  url += splits[4]
  urlYear += splits[4]
          //   and name to urlYear to get year
  urlYear += "/name"


          // Gets the url for species
  responeSpecies, err := http.Get(url)
    if err != nil {
      fmt.Print(err.Error())
      os.Exit(1)
    }
            // reads in the api for species
  responseData, err := ioutil.ReadAll(responeSpecies.Body)
    if err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusForbidden)
    }
            // declaring a struct and filling inn the json and api for species
  var responseObject Species
  json.Unmarshal(responseData, &responseObject)


          // Gets the url for year
  responeYear, err := http.Get(urlYear)
    if err != nil {
      panic(err)
    }
          // reads in the api for year
  responseYearData, err := ioutil.ReadAll(responeYear.Body)
    if err != nil {
      panic(err)
    }
          // filling inn the json and api for year
  json.Unmarshal(responseYearData, &responseObject)


          // Declaring and converting the struct into bytes
  speciesBytes, err := json.Marshal(responseObject)
          // writing out the bytes and status
  w.Header().Set("Content-Type", "application/json")
  w.Write(speciesBytes)
  w.WriteHeader(http.StatusOK)
}




      // Returns application status
func diag(w http.ResponseWriter, r *http.Request){
          // Gets the information about the website
  gbif, err := http.Get("http://api.gbif.org/v1/species/")
  if err != nil {
    panic(err)
  }
  restCountries, err := http.Get("http://restcountries.eu/rest/v2/")
  if err != nil {
    panic(err)
  }

          // Declaring the struct
  var responseObject Diagnostics

          // Filling in the struct
  responseObject.GBIF = gbif.StatusCode
  responseObject.RestCountries = restCountries.StatusCode
  responseObject.Version = "v1"
  responseObject.UpTime = int(time.Since(StartTime).Seconds())

          // Declaring and converting the struct into bytes
  diagBytes, err := json.Marshal(responseObject)
          // writing out the bytes and status
  w.Header().Set("Content-Type", "application/json")
  w.Write(diagBytes)
  w.WriteHeader(http.StatusOK)
}




      // Main functions
func main() {
          // Starts a timer at start of the program
  StartTime = time.Now().UTC()

          // Sets the localhost port to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

          // Functions retriving the url
          //  and starting the functions for the endpoints
	http.HandleFunc("/conservation/v1/country/", country)
  http.HandleFunc("/conservation/v1/species/", species)
  http.HandleFunc("/conservation/v1/diag/", diag)
          // Lets the user reenter input again and again
  http.ListenAndServe(":" + port, nil)
}
