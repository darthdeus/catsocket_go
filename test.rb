require "minitest/spec"
require "minitest/autorun"
require "redis"
require "restclient"
require "pry"
require "json"

redis = Redis.new

describe "Catsocket" do

  before :each do
    redis.flushdb
    redis.sadd("keys", "somekey")
  end

  it "saves the data in redis" do
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats", guid: "1" })

    response = redis.zrange "736f6d656b6579666f6fda39a3ee5e6b4b0d3255bfef95601890afd80709", 0, 99
    assert_equal ["1|cats"], response
  end

  it "works" do
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats", guid: "1" })
    response = RestClient.get("localhost:5000/?channel=foo&api_key=somekey&timestamp=1&guid=2")

    json = JSON.parse(response)

    assert_equal ["cats"], json["data"]
  end

  it "doesn't send back messages with the same GUID" do
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats", guid: "1" })
    response = RestClient.get("localhost:5000/?channel=foo&api_key=somekey&timestamp=1&guid=1")

    json = JSON.parse(response)

    assert_equal 0, json["data"].length, "Client shouldn't receive any data with his own GUID"
  end

  it "only receives messages with different guid" do
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "cats", guid: "1" })
    RestClient.post("localhost:5000", { channel: "foo", api_key: "somekey", data: "lolcats", guid: "2" })

    response = RestClient.get("localhost:5000/?channel=foo&api_key=somekey&timestamp=1&guid=1")

    json = JSON.parse(response)

    assert_equal ["lolcats"], json["data"]
  end


end
