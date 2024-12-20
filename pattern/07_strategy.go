/*
Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Паттерн «стратегия» позволяет заменять алгоритмы выполнения опредёленного действия на лету.

+ гибкость: можно легко добавлять новые алгоритмы
+ читаемость: заменяем множество условных операторов на цепочку алгоритмов
+ функциональность: можно выбирать алгоритм во время выполнения
- сложность: приводит к увеличению количества структур и интерфейсов


Можно использовать в системе платежей. Различные способы оплаты: кредитные карты, банковский перевод, крипта, paypal.
Маршрутизация - выбор маршрута, в зависимости от условий: пробки, ремонт дороги.
*/

package pattern

import (
	"fmt"
	"time"
)

type Strategy interface {
	Evict(c *Cache)
}

type TtlStrategy struct{}

func (l *TtlStrategy) Evict(c *Cache) {
	// ...
	fmt.Println("Evicting by TTL strategy")
}

type LruStrategy struct{}

func (l *LruStrategy) Evict(c *Cache) {
	// ...
	fmt.Println("Evicting by LRU strategy")
}

type LfuStrategy struct{}

func (l *LfuStrategy) Evict(c *Cache) {
	// ...
	fmt.Println("Evicting by LFU strategy")
}

type Cache struct {
	storage       map[string]*Elem
	evictStrategy Strategy
	size          int
	capacity      int
}

type Elem struct {
	data   string
	freq   int
	ttl    time.Time
	usedAt time.Time
}

func NewCashe() *Cache {
	return &Cache{
		storage:  make(map[string]*Elem),
		size:     0,
		capacity: 2,
	}
}

func (c *Cache) SetStrategy(e Strategy) {
	c.evictStrategy = e
}

func (c *Cache) Add(key, value string) {
	if c.size >= c.capacity {
		c.Evict()
	}
	c.size++
	c.storage[key] = &Elem{
		data:   value,
		ttl:    time.Now().Add(time.Hour),
		freq:   1,
		usedAt: time.Now(),
	}
}

func (c *Cache) Get(key string) (string, error) {
	if value, ok := c.storage[key]; ok {
		value.freq++
		value.usedAt = time.Now()
		value.ttl = time.Now().Add(time.Hour)
		return value.data, nil
	}
	return "", fmt.Errorf("key not found")
}

func (c *Cache) Evict() {
	c.evictStrategy.Evict(c)
	c.size--
}

func mainStrategy() {
	ttl := &TtlStrategy{}
	lru := &LruStrategy{}
	lfu := &LfuStrategy{}
	cache := NewCashe()
	cache.SetStrategy(ttl)
	cache.Add("1", "1")
	cache.Add("2", "2")
	cache.Add("3", "3")
	cache.SetStrategy(lru)
	cache.Add("4", "4")
	cache.Add("5", "5")
	cache.Add("6", "6")
	cache.SetStrategy(lfu)
	cache.Add("7", "7")
	cache.Add("8", "8")
	cache.Add("9", "9")
	cache.Add("10", "10")
}
