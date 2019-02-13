package metrics

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testTags = Tags{
		`tag0`:           0,
		`tag1`:           1,
		`success`:        true,
		`hello`:          `world`,
		`service`:        `rotator`,
		`server`:         `idk`,
		`worker_id`:      -1,
		`defaultTagBool`: true,
	}
)

func initDefaultTags() {
	defaultTags = *Tags{
		`defaultTag0`:       0,
		`defaultTagString`:  `string`,
		`defaultTagBool`:    false,
		`defaultOneMoreTag`: nil,
	}.ToFastTags()
}

func BenchmarkList(b *testing.B) {
	initDefaultTags()
	tags := Tags{
		`tag0`:       0,
		`tagString`:  `string`,
		`tagBool`:    false,
		`oneMoreTag`: nil,
	}
	for i := 0; i < 10000; i++ {
		tags[`value`] = i
		GaugeInt64(`test_metric`, tags)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		List()
	}
}

func BenchmarkGenerateStorageKey(b *testing.B) {
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := generateStorageKey(TypeCount, `test`, nil)
			if buf != nil {
				buf.Release()
			}
		}
	})
}

func BenchmarkGet(b *testing.B) {
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Get(TypeCount, `test`, nil)
		}
	})
}

func BenchmarkRegistry(b *testing.B) {
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GaugeInt64(``, nil)
		}
	})
}

func BenchmarkRegistryReal(b *testing.B) {
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GaugeInt64(`test_key`, testTags)
		}
	})
}
func BenchmarkAddToRegistryReal(b *testing.B) {
	var i uint64
	testTags[`i`] = &i
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.StoreUint64(&i, 1)
			GaugeInt64(`test_key`, testTags)
		}
	})
}

func BenchmarkRegistryRealReal(b *testing.B) {
	SetHiddenTags(HiddenTags{HiddenTag{`success`, nil}, HiddenTag{`campaign_id`, ExceptValues{1}}})
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testTags := Tags{
				`tag0`:           0,
				`tag1`:           1,
				`success`:        true,
				`hello`:          `world`,
				`service`:        `rotator`,
				`server`:         `idk`,
				`worker_id`:      -1,
				`defaultTagBool`: true,
			}
			GaugeInt64(`test_key`, testTags)
		}
	})
	SetHiddenTags(nil)
}

/*
func BenchmarkRegistryReal_FastTags_withHiddenTag(b *testing.B) {
	SetHiddenTags(HiddenTags{HiddenTag{`success`, nil}, HiddenTag{`campaign_id`, ExceptValues{1}}})
	initDefaultTags()
	testTagsFast := testTags.ToFastTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GaugeInt64(`test_key`, testTagsFast)
		}
	})
	SetHiddenTags(nil)
}

func BenchmarkRegistryReal_FastTags(b *testing.B) {
	initDefaultTags()
	testTagsFast := testTags.ToFastTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GaugeInt64(`test_key`, testTagsFast)
		}
	})
}*/

func BenchmarkTagsString(b *testing.B) {
	initDefaultTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := generateStorageKey(TypeGaugeInt64, `testKey`, testTags)
			buf.Release()
		}
	})
}

func BenchmarkTagsFastString(b *testing.B) {
	initDefaultTags()
	testTagsFast := testTags.ToFastTags()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := generateStorageKey(TypeGaugeInt64, `testKey`, testTagsFast)
			buf.Release()
		}
	})
}

func TestGC(t *testing.T) {
	var memstats, cleanedMemstats runtime.MemStats
	goroutinesCount := runtime.NumGoroutine()
	runtime.GC()
	runtime.ReadMemStats(&memstats)
	metric := GaugeInt64(`test_metric`, nil)
	//newGoroutinesCount := runtime.NumGoroutine()
	metric.Stop()
	GC()
	runtime.GC()
	runtime.ReadMemStats(&cleanedMemstats)
	cleanedGoroutinesCount := runtime.NumGoroutine()
	//assert.Equal(t, goroutinesCount+1, newGoroutinesCount)
	assert.Equal(t, cleanedGoroutinesCount, goroutinesCount)
	//assert.Equal(t, memstats.HeapInuse, cleanedMemstats.HeapInuse)
}

func TestRegistry(t *testing.T) {
	defaultTags = *Tags{
		`datacenter`: `EU`,
		`hostcode`:   `999`,
		`hostname`:   `e0df6242fcbf`,
		`service`:    `rotator`,
	}.ToFastTags()

	tags := Tags{
		`code`:      400,
		`format_id`: `unknown`,
		`network`:   `unknown`,
	}

	tags0 := tags.Copy()
	tags0[`key`] = `dsp.bid`
	GaugeInt64(`requests`, tags0)

	tags1 := tags.Copy()
	tags1[`key`] = `dsp.bid.tjnative`
	GaugeInt64(`requests`, tags1)

	assert.Equal(t, `dsp.bid`, Get(TypeGaugeInt64, `requests`, tags0).GetTag(`key`))
	assert.Equal(t, `dsp.bid.tjnative`, Get(TypeGaugeInt64, `requests`, tags1).GetTag(`key`))
}

func TestTagsString(t *testing.T) {
	initDefaultTags()
	{
		buf := generateStorageKey(TypeGaugeInt64, `testKey`, testTags)
		assert.Equal(t, `testKey,defaultOneMoreTag=null,defaultTag0=0,defaultTagBool=false,defaultTagString=string,hello=world,server=idk,service=rotator,success=true,tag0=0,tag1=1,worker_id=-1@gauge_int64`, buf.result.String())
		buf.Release()
	}

	{
		SetHiddenTags(HiddenTags{HiddenTag{`app_id`, nil}, HiddenTag{`spot`, nil}, HiddenTag{`spot_id`, nil}, HiddenTag{`app`, nil}, HiddenTag{`campaign_id`, ExceptValues{123}}, HiddenTag{`user_id`, ExceptValues{12}}})
		assert.Equal(t, `app`, metricsRegistry.getHiddenTags()[0].Key)

		tags := Tags{`spot`: true, `campaign_id`: 123, `user_id`: 55}
		considerHiddenTags(tags)
		buf := generateStorageKey(TypeGaugeInt64, `testKey`, tags)
		assert.Equal(t, `testKey,defaultOneMoreTag=null,defaultTag0=0,defaultTagBool=false,defaultTagString=string,campaign_id=123,spot=hidden,user_id=hidden@gauge_int64`, buf.result.String())
		buf.Release()
	}
}