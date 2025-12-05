import {useState} from 'react'
import thyngoLogo from '/thyngo-dark.svg'
import '../styles/App.css'
import GetPosts from "../components/GetPosts.tsx";

function App() {
    const [count, setCount] = useState(0)

    return (
        <>
            <div>
                <a href="https://raph.mailhiotinc.com" target="_blank">
                    <img src={thyngoLogo} className="logo" alt="thynGo logo"/>
                </a>
            </div>
            <h1>Dashboard pour <span style={{fontWeight: "300", fontSize: "4rem"}}>ThynGo</span></h1>

            <div className="card">
                <button onClick={() => setCount((count) => count + 1)}>
                    count is {count}
                </button>
            </div>

            <GetPosts/>

        </>
    )
}

export default App
