declare interface Subscribe {
    id: number;
    created_at: string;
    updated_at: string;
    url: string;
    alias: string;
    nodes: NodeInfo[];
}

declare interface NodeInfo {
    id: number;
    created_at: string;
    updated_at: string;

    subscription_id: number;
    node_type: string;
    remarks: string;
    uuid: string;
    address: string;
    port: number;
    ns: string;
    tcping: number;
    ascription: string;
    enabled: boolean;
    //vmess参数
    v: string;
    scy: string;//加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero
    aid: number;
    //vless参数
    flow: string;//流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
    encryption: string;

    network: string; ////传输协议 tcp,kcp,ws,h2,quic,grpc
    type: string;  //伪装类型 ws,h2,grpc：无    tcp：none，http    kcp,quic：none，srtp，utp，wechat-video，dtls，wireguard
    host: string;
    path: string;
    mode: string; //grpc传输模式 gun，multi

    security: string;        //传输层安全类型 none,tls,reality
    sni: string;             //
    fp: string;              //

    alpn: string;            //tls
    allowInsecure: boolean;  //tls

    pbk: string;             //reality
    sid: string;             //reality
    spx: string;             //reality
}

//配置参数
declare interface Config {
    id: number;
    created_at: string;
    updated_at: string;
    domestic_type: string;
    abroad_type: string;
    auto_change_node: string;
    node_pool_model:string;
    host: string;
    wifi_proxy: string;
    ipv6_net: string;
    allow_outside_tcp_udp: string;
    allow_apps: [];
}
