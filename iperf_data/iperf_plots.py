import matplotlib.pyplot as plt

def plot_name_dual(name):
    none_file = name + "none.iperf.log"
    tcp_file = name + "tcp.iperf.log"
    rrtcp_file = name + "rrtcp.iperf.log"

    with open(none_file) as f:
        none_lines = f.readlines()
    with open(tcp_file) as f:
        tcp_lines = f.readlines()
    with open(rrtcp_file) as f:
        rrtcp_lines = f.readlines()

    times = range(60)
    times = [time/2.0 for time in times]

    id1 = none_lines[0].split(",")[-4]
    id1_str = "," + id1 + "," 

    none_speed1 = []
    none_speed2 = []
    for line in none_lines[:-2]:
        speed = int(line.split(",")[-1])
        if id1_str in line: 
            none_speed1.append(speed)
        else:
            none_speed2.append(speed)
     
    tcp_speed1 = []
    tcp_speed2 = []
    for line in tcp_lines[:-2]:
        speed = int(line.split(",")[-1])
        if id1_str in line: 
            tcp_speed1.append(speed)
        else:
            tcp_speed2.append(speed)
     
    rrtcp_speed1 = []
    rrtcp_speed2 = []
    for line in rrtcp_lines[:-2]:
        speed = int(line.split(",")[-1])
        if id1_str in line: 
            rrtcp_speed1.append(speed)
        else:
            rrtcp_speed2.append(speed)

    none_speed1 = none_speed1[:len(times)]
    none_speed2 = none_speed2[:len(times)]
    tcp_speed1 = tcp_speed1[:len(times)]
    tcp_speed2 = tcp_speed2[:len(times)]
    rrtcp_speed1 = rrtcp_speed1[:len(times)]
    rrtcp_speed2 = rrtcp_speed2[:len(times)]

    plt.title(name + " Throughput Over Time")
    plt.ylabel("Bits of Throughput")
    plt.xlabel("Seconds")
    plt.plot(times, none_speed1, color='b', linestyle='--', label='None A')
    plt.plot(times, none_speed2, color='b', label='None B')
    plt.plot(times, tcp_speed1, color='g', linestyle='--', label='TCP A')
    plt.plot(times, tcp_speed2, color='g', label='TCP B')
    plt.plot(times, rrtcp_speed1, color='r', linestyle='--', label='RRTCP A')
    plt.plot(times, rrtcp_speed2, color='r', label='RRTCP B')
    plt.legend()
    plt.savefig('../plots/iperf/' + name + '.png')
    plt.savefig('../plots/iperf/' + name + '.eps')
    plt.show()

def plot_name_single(name):
    none_file = name + "none.iperf.log"
    tcp_file = name + "tcp.iperf.log"
    rrtcp_file = name + "rrtcp.iperf.log"

    with open(none_file) as f:
        none_lines = f.readlines()
    with open(tcp_file) as f:
        tcp_lines = f.readlines()
    with open(rrtcp_file) as f:
        rrtcp_lines = f.readlines()

    times = range(30)
    #times = [time/2.0 for time in times]

    none_speed = []
    for line in none_lines[:-1]:
        speed = int(line.split(",")[-1])
        none_speed.append(speed)
     
    tcp_speed = []
    for line in tcp_lines[:-1]:
        speed = int(line.split(",")[-1])
        tcp_speed.append(speed)

    rrtcp_speed = []
    for line in rrtcp_lines[:-1]:
        speed = int(line.split(",")[-1])
        rrtcp_speed.append(speed)

    plt.title(name + " Throughput Over Time")
    plt.ylabel("Bits of Throughput")
    plt.xlabel("Seconds")
    plt.plot(times, none_speed, color='b', label='None')
    plt.plot(times, tcp_speed, color='g', label='TCP')
    plt.plot(times, rrtcp_speed, color='r', label='RRTCP')
    plt.legend()
    plt.savefig('../plots/iperf/' + name + '.png')
    plt.savefig('../plots/iperf/' + name + '.eps')
    plt.show()

plot_name_single('estonia-')
