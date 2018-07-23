#include <stdlib.h>
#include <sodium.h>
#include <inttypes.h>
#include <unistd.h>
#include <sys/time.h>

static unsigned char buf[1024*1024];

typedef struct rando {
	uint64_t v0, v1, v2, v3, v4;
	uint64_t id, cnt, w0, w1;
	uint32_t pos;
} rando;

static uint64_t seedspace[128];

static void permute(rando* s) {
	uint64_t v0, v1, v2, v3;
	uint32_t pos;
	v0 = s->v0;
	v1 = s->v1;
	v2 = s->v2;
	v3 = s->v3;
	pos = s->pos;
	s->pos = s->pos*9+1;
	v0 ^= seedspace[pos>>25];
	v1 ^= s->id;
	v3 ^= s->cnt;
	s->cnt += s->pos;
	for (int i = 0; i < 5; i++) {
		v0 += v1;
		v1 = v1<<13 | v1>>(64-13);
		v1 ^= v0;
		v0 = v0<<32 | v0>>(64-32);

		v2 += v3;
		v3 = v3<<16 | v3>>(64-16);
		v3 ^= v2;

		v0 += v3;
		v3 = v3<<21 | v3>>(64-21);
		v3 ^= v0;

		v2 += v1;
		v1 = v1<<17 | v1>>(64-17);
		v1 ^= v2;
		v2 = v2<<32 | v2>>(64-32);
	}
	s->v0 = v0;
	s->v1 = v1;
	s->v2 = v2;
	s->v3 = v3;
}

double nowtime() {
	struct timeval now;
	gettimeofday(&now, NULL);
	return (double)now.tv_sec + (double)now.tv_usec / 1e6;
}

int main(void) {
	rando s;
	double prev, now;
	if (sodium_init() < 0) {
		abort();
	}
	randombytes_buf(&seedspace, sizeof(seedspace));
	randombytes_buf(&s, sizeof(s));
	prev = nowtime();
	while(1) {
		int i, r;
		for (i=0; i<sizeof(buf); i+=16) {
			permute(&s);
			*(uint64_t*)(buf+i) = s.v1 ^ s.w0;
			*(uint64_t*)(buf+i+8) = s.v3 ^ s.w1;
		}
		r = write(1, buf, sizeof(buf));
		if (r < 0)
			break;
		now = nowtime();
		if (now - prev > 0.1) {
			randombytes_buf(&seedspace, sizeof(seedspace));
			prev = now;
		}
	}
}
