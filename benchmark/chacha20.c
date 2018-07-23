#include <stdlib.h>
#include <sodium.h>
#include <inttypes.h>
#include <unistd.h>

static unsigned char buf[1024*1024];
int main(void) {
	unsigned char key[crypto_stream_chacha20_KEYBYTES];
	// nonce is 8 byte
	uint64_t nonce = 0;
	ssize_t r;
	if (sodium_init() < 0) {
		abort();
	}
	randombytes_buf(key, sizeof(key));
	while(1) {
		if (crypto_stream_chacha20(buf, sizeof(buf), (unsigned char*)&nonce, key) < 0) {
			abort();
		}
		r = write(1, buf, sizeof(buf));
		if (r < 0)
			break;
		nonce++;
	}
}
