import React, {useEffect, useState, useRef } from 'react'
import Select from "react-select";
import {Countries} from "../../constants/Countries";
import {useHistory} from "react-router-dom";
import ErrorBox from "../ErrorBox";
import {getAxiosInstance} from "../Auth";
import {ConvertImageToString} from "../utilities/images";

export default function Settings() {

    const history = useHistory();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [country, setCountry] = useState("");
    const [defaultCountry, setDefaultCountry] = useState({});
    const [avatar, setAvatar] = useState(null);
    const [error, setError] = useState("")


    useEffect(() => {
        const getUser = async () => {
            const instance = getAxiosInstance();
            const data = await instance.get("/api/v1/user/" + localStorage.getItem("email"));
            const user = data.data;

            user.avatarFilename = "/api/v1/image/" + user.Avatar;
            setEmail(user.Email);
            setName(user.Name);
            setCountry(user.Country);
            const c = Countries.find(obj => obj.value === user.Country);
            setDefaultCountry(c);
        };

        const fetchData = async() => {
            await getUser();
        };
        fetchData();
    }, []);


    const HandleFormSubmit = async (e) => {
        e.preventDefault();

        const data = { email, password, country, name, avatar };

        if(email === "") {
            setError("Email address is required");
            return false;
        }
        if(name === "") {
            setError("Name is required");
            return false;
        }
        if(country === "") {
            setError("Country is required")
            return false;
        }

        if(password === "") {
            delete data.password;
        }

        if(avatar !== null) {
            data.avatar = await ConvertImageToString(avatar);
        }

        const instance = getAxiosInstance();

        const updateUserUrl = "/api/v1/user";
        await instance.patch(
            updateUserUrl,
            data
        ).then((data) => {
            // Successful login
            localStorage.setItem("email", email);
            localStorage.setItem("name", name);

            window.location.reload();

        }).catch((error) => {
            console.log(error);
            setError(error.response.data.error);
            return false;
        });

    }

    const handleCountry = (e) => {
        setCountry(e.value);
        setDefaultCountry(e);
    }

    const handleAvatar = (e) => {
        setAvatar(e.target.files[0]);
    }

    // setCountry("US");

    return (
        <div className="settings-form form-540">
            <h2>Settings</h2>
            <form onSubmit={HandleFormSubmit}>
                <label>
                    Email address*
                    <input type="text" name="email" placeholder="Enter email" value={email} onChange={e => setEmail(e.target.value)}/>
                </label>
                <label>
                    Password
                    <input type="password" name="password" placeholder="Password"
                           onChange={e => setPassword(e.target.value)}/>
                </label>
                <label>
                    Name*
                    <input type="text" name="name" placeholder="Your name (shown publicly)" value={name}
                           onChange={e => setName(e.target.value)}/>
                </label>
                <label>
                    Country*
                    <Select name="country" value={defaultCountry} options={Countries} onChange={handleCountry}/>
                </label>
                <label>
                    Avatar
                    <div>
                        <input type="file" name="avatar" onChange={handleAvatar}/>
                    </div>
                </label>
                <div>
                    <input type="submit" value="Save" className="orange-button raise"/>
                </div>

                {error &&
                <ErrorBox message={error}/>
                }
            </form>
        </div>
    );
}