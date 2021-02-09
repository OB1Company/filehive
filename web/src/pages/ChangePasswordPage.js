import React, {useEffect, useState} from 'react'
import {Helmet} from "react-helmet";
import Header from "../Header";
import Footer from "../Footer";
import {Link} from "react-router-dom";
import ErrorBox, {SuccessBox} from "../components/ErrorBox";
import {getAxiosInstance} from "../components/Auth";
import {UseQuery} from "../components/utilities/images";


export default function ChangePasswordPage() {

    let query = UseQuery();

    const code = query.get("code");
    const email = query.get("email");

    useEffect(()=>{
        // Check if valid email/code is valid
        const submitRequest = async ()=>{
            const instance = getAxiosInstance();
            instance.get("/api/v1/checkresetcode?email="+email+"&code="+code)
                .then((response)=>{
                    console.log(response);
                    if(!response.data.success) {
                        window.location.href='/';
                    }
                })
                .catch((err)=>{
                    console.log(err);
                })
        }
        submitRequest();
    }, []);

    const [password, setPassword] = useState("");
    const [confirm, setConfirm] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const HandlePasswordReset = (e)=>{
        e.preventDefault();
        setError("");

        if(password !== confirm) {
            setError("passwords must match");
            return
        }

        const data = {
            email: email,
            password: password,
            code: code,
        };

        const submitRequest = async ()=>{

            const resetSuccess = () => {
                return (
                    <div>Your password has been updated. <a href="/login">Login</a></div>
                );
            }

            const instance = getAxiosInstance();
            instance.post("/api/v1/passwordreset",
                data)
                .then((data) => {
                    setSuccess(resetSuccess);
                    setError("");
                })
                .catch((e) => {
                    setError(e.response.data.error);
                    setSuccess("");
                    return false;
                });
        }
        submitRequest();

    }

    return (
        <div className="container">
            <Helmet>
                <title>Filehive | Change Password</title>
            </Helmet>
            <Header/>
            <div className="subBody">
                <div className="maincontent">

                    <div className="Login form-540">
                        <h2>Change Password</h2>

                        <form onSubmit={HandlePasswordReset}>
                            <label>
                                New password
                                <input type="password" name="password" onChange={e => setPassword(e.target.value)}/>
                            </label>
                            <label>
                                Confirm password
                                <input type="password" name="confirm" onChange={e => setConfirm(e.target.value)}/>
                            </label>
                            <div>
                                <input type="submit" value="Reset password" className="raise orange-button" />
                                <Link to='/login'>Log in</Link>
                            </div>

                            {error &&
                            <ErrorBox message={error}/>
                            }
                            {success &&
                            <SuccessBox message={success}/>
                            }
                        </form>

                    </div>

                </div>
            </div>
            <Footer/>
        </div>
    )
}