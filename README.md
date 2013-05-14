Go-Redisbayes
=============

RedisBayes Golang Version


What Is This?
=============

It's a spam filter.  I wrote this to filter spammy comments from a high
traffic forum website and it worked pretty well.  It can work for you too :)
It's not tied to any particular format like email, it just deals with the raw
text.


Installation
============

go get github.com/pengfei-xue/redisbayes


Basic Usage
===========

::

    package main
    
    import (
    	"fmt"
    	"github.com/pengfei-xue/redisbayes"
    )
    
    // you should add a config file config.yaml in your working dir
    // this will be changed to use the default config if not provided
    func main() {
    	redisbayes.Train("good", "sunshine drugs love sex lobster sloth")
    	redisbayes.Train("bad", "fear death horror government zombie god")
    
    	class := redisbayes.Classify("sloths are so cute i love them")
    	fmt.Println(class)
    	class = redisbayes.Classify("i fear god and love the government")
    	fmt.Println(class)
    
    	redisbayes.Untrain("good", "sunshine drugs love sex lobster sloth")
    	redisbayes.Untrain("bad", "fear death horror government zombie god")
    }
