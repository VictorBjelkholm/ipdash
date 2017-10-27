'use strict'

const React = require('react')
const IPFS = require('ipfs')

const PieChart = require("react-chartjs").Pie

class App extends React.Component {
  componentDidMount () {
    const node = new IPFS({
      repo: String(Math.random() + Date.now()),
      EXPERIMENTAL: { // enable experimental features
        pubsub: true,
        sharding: true, // enable dir sharding
        wrtcLinuxWindows: true, // use unstable wrtc module on Linux or Windows with Node.js,
        dht: true // enable KadDHT, currently not interopable with go-ipfs
      }
    })

    node.on('ready', () => {
      setInterval(() => {
        node.swarm.peers((err, peers) => {
          if (err) throw err
          console.log("Connected peers:", peers.length)
        })
      }, 1000 * 5)
      setTimeout(() => {
        node.swarm.connect('/ip4/127.0.0.1/tcp/9999/ws/ipfs/QmcKrQreYe9AxebTBSvnDo6vem9DBiqsLB1d8UeXhkxuBr', (err) => {
          if (err) throw err
        })
      }, 1000)
      node.pubsub.subscribe('ipfs-dashboard', (msg) => {
        const stats = JSON.parse(msg.data.toString())
        console.log(stats)
      })
    })
  }

  render() {
    return <div>
      <div>
        <div>Global Upload Rate: 5.3 GB/s</div>
        <div>Global Download Rate: 5.3 GB/s</div>
        <div>Available Space: 5.3 GB/s</div>
        <div>Used Space: 5.3 GB/s</div>
      </div>
      <div>
        Map of peers
      </div>
      <div>
        <div>
          <PieChart data={}/>
        </div>
        <div>Piechart OS/Arch</div>
        <div>Number of CPUs</div>
      </div>
      <div>
        Table of connected peers
      </div>
    </div>
  }
}
module.exports = App
