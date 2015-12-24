with open('tcp_160_5.delay.log', 'r') as f:
    tcp_lines = f.readlines() 
with open('rrtcp_160_5.delay.log', 'r') as f:
    rrtcp_lines = f.readlines() 

tcp_send_avg = 0
tcp_recv_avg = 0
tcp_wire_avg = 0
for line in tcp_lines:
    split_line = line.split()
    app_send = int(split_line[0])
    net_send = int(split_line[1])
    net_recv = int(split_line[2])
    app_recv = int(split_line[3])
    
    if app_send > net_send or net_send > net_recv or net_recv > app_recv:
        print line
        raw_input()

    tcp_send_avg += (net_send - app_send)
    tcp_recv_avg += (app_recv - net_recv)
    tcp_wire_avg += (net_recv - net_send)
    
tcp_send_avg = float(tcp_send_avg)/len(tcp_lines)
tcp_recv_avg = float(tcp_recv_avg)/len(tcp_lines)
tcp_wire_avg = float(tcp_wire_avg)/len(tcp_lines)

print "TCP Send Avg: ", tcp_send_avg
print "TCP Recv Avg: ", tcp_recv_avg
print "TCP Wire Avg: ", tcp_wire_avg

rrtcp_send_avg = 0
rrtcp_recv_avg = 0
rrtcp_wire_avg = 0
for line in rrtcp_lines:
    split_line = line.split()
    app_send = int(split_line[0])
    net_send = int(split_line[1])
    net_recv = int(split_line[2])
    app_recv = int(split_line[3])
    
    if app_send > net_send or net_send > net_recv or net_recv > app_recv:
        print line
        raw_input()

    rrtcp_send_avg += (net_send - app_send)
    rrtcp_recv_avg += (app_recv - net_recv)
    rrtcp_wire_avg += (net_recv - net_send)
    
rrtcp_send_avg = float(rrtcp_send_avg)/len(rrtcp_lines)
rrtcp_recv_avg = float(rrtcp_recv_avg)/len(rrtcp_lines)
rrtcp_wire_avg = float(rrtcp_wire_avg)/len(rrtcp_lines)

print "RRTCP Send Avg: ", rrtcp_send_avg
print "RRTCP Recv Avg: ", rrtcp_recv_avg
print "RRTCP Wire Avg: ", rrtcp_wire_avg
