import React from "react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { apiService } from "../axios";
import "../components/login.css"
import logo from "../assets/logo.png"
import username from "../assets/username.svg";
import password from "../assets/password.svg";

const Login = () => {
    const {
        register,
        handleSubmit,
        formState: { errors, isValid },
    } = useForm({ mode: "onChange" });
    const [errorMessage, setErrorMessage] = useState("");
    const navigate = useNavigate();

    const onSubmit = async (data) => {
        
        try {
            const response = await apiService.login({
                Username: data.username,
                Password: data.password,
            })
            sessionStorage.setItem("token", `${response.data.Data.token}`);
            navigate("/home");
        } catch (error) {
            switch (error.response.data.error) {
                case "client not found":
                    setErrorMessage("Cliente não cadastrado");
                    break;
                case "invalid credentials":
                    setErrorMessage("Credenciais inválidas");
                    break;
                case "more than one user logged":
                    setErrorMessage("Um dispositivo já está conectado com a conta.")
                    break;
            }   
        }        
    };
    return (
        <div className="login">


            <form onSubmit={handleSubmit(onSubmit)} className="login-container">
                <img src={logo} alt="" />
                <div className="username-group">
                    <img src={username} alt="" width={"30px"} />
                    <input
                        className="login-input"
                        type="text"
                        {...register("username", { required: true })}
                        placeholder="Nome de usuário"
                    />
                </div>
                    {errors.username && <h5 className="advice">Insira o username.</h5>}

                <div className="password-group">
                    <img src={password} alt="" width={"30px"}/>
                    <input
                        className="login-input"
                        type="password"
                        name="password"
                        {...register("password", { required: true })}
                        placeholder="Senha"
                    />
                </div>
                {errors.password && <h5 className="advice">Insira a senha.</h5>}

                <button type="submit" className={`button ${!isValid ? "disabled" : ""}`} disabled={!isValid}>
                    Entrar
                </button>

                {errorMessage && <div className="advice">{errorMessage}</div>}
            </form>
        </div>
    );
};

export default Login;