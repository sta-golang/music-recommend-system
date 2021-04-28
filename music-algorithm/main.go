package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/music-algorithm/music2vec"
	"github.com/ynqa/wego/pkg/embedding"
	"github.com/ynqa/wego/pkg/model/modelutil/vector"
	"github.com/ynqa/wego/pkg/model/word2vec"
	"github.com/ynqa/wego/pkg/search"
)

func main() {
	//word2()
	//readModel()
	music2()
	//testSort()
}

func music2() {
	model, _ := music2vec.NewDefaultModel()
	bys, err := ioutil.ReadFile("songVec.txt")
	if err != nil {
		panic(err)
	}
	err = model.LoadData(bys)
	if err != nil {
		panic(err)
	}
	err = model.Train()
	if err != nil {
		panic(err)
	}
	fmt.Println(model.SimilarMusic(296542).TopN(10))
}

func word2() {

	model, err := word2vec.New(
		word2vec.Window(5),
		word2vec.Model(word2vec.Cbow),
		word2vec.Optimizer(word2vec.NegativeSampling),
		word2vec.NegativeSampleSize(3),
		word2vec.Verbose(),
		word2vec.Dim(10),
		word2vec.Iter(200),
	)

	if err != nil {
		// failed to create word2vec.
	}
	bys, err := ioutil.ReadFile("kkk.txt")
	if err != nil {
		panic(err)
	}
	newReader := strings.NewReader(str.BytesToString(bys))

	if err = model.Train(newReader); err != nil {
		// failed to train.
		fmt.Println("err = ", err)
	}

	model.Save(os.Stdin, vector.Single)
	// write word vector.
	fmt.Println(model.WordVector(vector.Single).Row())
	fmt.Println(model.WordVector(vector.Single).Col())
	fmt.Println(model.WordVector(vector.Single).Slice(0))
	fmt.Println(model.WordVector(vector.Single).Slice(1))
}

func readModel() {
	model, err := word2vec.New(
		word2vec.Window(5),
		word2vec.Model(word2vec.Cbow),
		word2vec.Optimizer(word2vec.NegativeSampling),
		word2vec.NegativeSampleSize(3),
		word2vec.Verbose(),
		word2vec.Dim(10),
		word2vec.Iter(200),
	)
	if err != nil {
		panic(err)
	}
	file, err := os.Open("helo.txt")
	if err != nil {
		panic(err)
	}
	if err = model.Train(file); err != nil {
		panic(err)
	}
	model.Save(os.Stdin, vector.Single)

}

func query() {
	input, err := os.Open("word.txt")
	if err != nil {
		panic(err)
	}
	defer input.Close()
	embs, err := embedding.Load(input)
	if err != nil {
		panic(err)
	}
	searcher, err := search.New(embs...)
	if err != nil {
		panic(err)
	}
	neighbors, err := searcher.SearchInternal("given_word", 10)
	if err != nil {
		panic(err)
	}
	neighbors.Describe()
}
