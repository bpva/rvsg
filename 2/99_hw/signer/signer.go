package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)

		go func(j job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			j(in, out)
		}(j, in, out, wg)

		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	w := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {
		dataStr := fmt.Sprintf("%v", data)
		w.Add(1)

		go func(dataStr string) {
			defer w.Done()

			mu.Lock()
			md5Data := DataSignerMd5(dataStr)
			mu.Unlock()

			crc32DataChan := make(chan string)
			crc32Md5Chan := make(chan string)

			go func() {
				crc32DataChan <- DataSignerCrc32(dataStr)
			}()

			go func() {
				crc32Md5Chan <- DataSignerCrc32(md5Data)
			}()

			crc32Data := <-crc32DataChan
			crc32Md5 := <-crc32Md5Chan

			out <- crc32Data + "~" + crc32Md5
		}(dataStr)
	}

	w.Wait()
}

func MultiHash(in, out chan interface{}) {
	w := &sync.WaitGroup{}

	for data := range in {
		dataStr := fmt.Sprintf("%v", data)
		w.Add(1)

		go func(dataStr string) {
			defer w.Done()

			results := make([]string, 6)
			var innerWg sync.WaitGroup

			for th := 0; th <= 5; th++ {
				innerWg.Add(1)

				go func(th int) {
					defer innerWg.Done()
					data := fmt.Sprintf("%d", th) + dataStr
					results[th] = DataSignerCrc32(data)
				}(th)
			}

			innerWg.Wait()
			out <- strings.Join(results, "")
		}(dataStr)
	}

	w.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}
