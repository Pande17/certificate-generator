import { useEffect, useState } from "react"
import axios from "axios"
import{message} from antd


const competence = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])

  useEffect(() => {
    const fetchingData = async() => {
      setLoading(true)
      try{
        const response = await axios.get(
          `http://127.0.0.1:3000/api/competence`
        );
        const data = response.data.data
      }catch(err){
        console.error('error : ', err)
        message.error('error : ', err)
      }finally {
         setLoading(false)
      }
    }
  })

}