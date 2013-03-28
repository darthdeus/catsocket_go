require "minitest/spec"
require "minitest/autorun"
require "redis"
require "restclient"
require "pry"
require "json"

redis = Redis.new

describe "Catsocket" do

  it "saves the data in redis" do
    redis.sadd("keys", "somekey")
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats" })

    response = redis.zrange "736f6d656b6579666f6fda39a3ee5e6b4b0d3255bfef95601890afd80709", 0, 99
    assert response == ["cats"]
    redis.flushdb
  end

  it "works" do
    redis.sadd("keys", "somekey")
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats" })
    response = RestClient.get("localhost:5000/?channel=foo&api_key=somekey&timestamp=1")

    json = JSON.parse(response)

    assert json["data"] == ["cats"]
    redis.flushdb
  end

  it "doesn't send back messages with the same GUID" do
    redis.sadd("keys", "somekey")
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats", guid: "1" })
    response = RestClient.get("localhost:5000/?channel=foo&api_key=somekey&timestamp=1&guid=1")

    json = JSON.parse(response)

    assert json["data"].length == 0, "Client shouldn't receive any data with his own GUID"
    redis.flushdb
  end


end


