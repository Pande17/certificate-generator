import { useState,useEffect } from "react";
import axios from "axios";

const [field, setField] = useState([]);
const [selected, setSlected] = useState([]);
const [competence, setCompetence] = useState([]);


useEffect(()=>{
    const fetchData = async() =>{
        try{
            const response = await axios.get(
              "http://127.0.0.1:3000/api/competence"
            );
        }catch (err){
            console.log('error bg.Nih erorrnya : ', err)
        }
    };
    fetchData();
},[])

