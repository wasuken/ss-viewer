import React, { Component } from 'react';

class App extends Component{
    constructor(props){
        super(props);
        this.state = {
            json: [],
            displayContents: "",
            page: 0
        };
        this.getJson();
        this.clickAction = this.clickAction.bind(this);
        let that = this;
        fetch("/api/v1/ss/page").then(resp=>resp.json()).then(json =>{
            that.setState({
                json: [],
                displayContents: "",
                page: json.page
            });
        });
    }
    clickAction(){
        let result = this.state.json.find((map)=>{
            return map.Id == document.querySelector("#select").value;
        });
        this.setState({
            json: this.state.json,
            displayContents: result.Contents,
            page: this.state.page
        });
    }
    getJson(page=1){
        let that = this;
        fetch("/api/v1/ss?page=" + page).then((resp)=>resp.json()).then((json)=>{
            that.setState({
                json: json,
                displayContents: this.state.contents,
                page: this.state.page
            });
        });
    }
    getLink(){
        let lst = [];
        for(let i=1;i<this.state.page+1;i++){
            lst.push(<a className="page" href="#" onClick={() => this.getJson(i)}>{i}</a>);
        }
        return lst;
    }
    render(){
        let lst = this.getLink();
        
        let titles = this.state.json.map((map)=>(
            <option value={map.Id}>{map.Title}</option>
        ));
        return (
            <div>
              <div style={{float:"left", width:"30%"}}>
                <select id="select" onChange={this.clickAction}>
                  {titles}
                </select>
                <div>{lst}</div>
              </div>
              <div style={{float:"right", width:"70%", height: "100%",overflowY: "scroll"}}
                   className="text" dangerouslySetInnerHTML={{__html: this.state.displayContents}}>
              </div>
            </div>
        );
    }
}

export default App;
