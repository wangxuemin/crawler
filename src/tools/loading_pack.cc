#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <mc_pack.h>
#include <mc_pack_rp.h>
#include <nshead.h>

int
main(int argc, char **argv){
    if (argc != 2) {
        printf ("./loading_pack [file]\n");
        exit(-1);
    }
    bsl::ResourcePool rp;
    char buf[102400];
    int fd = open(argv[1], O_RDONLY);
    if (fd < 0) {
        perror("open");
        exit(-1);
    }
    memset(buf, 0, 102400);
    read(fd, buf, 102400);
    
    nshead_t *head = (nshead_t *)buf;
    if (head->magic_num != NSHEAD_MAGICNUM) {
        printf ("Invalid pack\n");
        return -1;
    }

    printf ("=====\nid:%d\nversion:%d\nlog_id:%d\nprovider:%s\nbody_len:%d\n=====\n", 
            head->id, head->version,
            head->log_id, head->provider,
            head->body_len);

    mc_pack_t* pack = mc_pack_open_r_rp(buf + sizeof(nshead_t), 102400, &rp);
    if (MC_PACK_PTR_ERR(pack) < 0) {
        printf ("open pack error!\n");
        exit(-1);
    }

    mc_pack_print(pack);

    return 0;
}
